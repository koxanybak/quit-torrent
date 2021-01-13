package download

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/koxanybak/quit-torrent/torrent/client"
	"github.com/koxanybak/quit-torrent/torrent/message"
	"github.com/koxanybak/quit-torrent/torrent/piece"
)

type pieceProgress struct {
	index			int
	client			*client.Client
	downloaded		int
	buf				[]byte
	requested		int
}

func (pp *pieceProgress) handleMessage(msg *message.Message) error {
	//fmt.Printf("Got message of type %s with payload length %d from peer %v\n", msg.ID.String(), len(msg.Payload), pp.client.Peer.IP)
	switch msg.ID {
	case message.MsgKeepAlive:
		return nil
	case message.MsgChoke:
		pp.client.Choked = true
	case message.MsgUnchoke:
		pp.client.Choked = false
	case message.MsgHave:
		pp.client.BitField.SetPiece(int(binary.BigEndian.Uint32(msg.Payload)))
	case message.MsgPiece:
		payload, err := piece.UnmarshalPiecePayload(msg.Payload)
		if err != nil {
			return err
		}

		// Save block to memory
		n := copy(pp.buf[payload.Begin : payload.Begin + len(payload.Block)], payload.Block)
		if n != len(payload.Block) {
			log.Panicf("Couldn't not save to memory block for piece %d of length %d beginning at %d from peer %v",
						payload.Index,
						len(payload.Block),
						payload.Begin,
						pp.client.Peer.IP)
		}

		pp.downloaded += n

		// fmt.Println("------------------------------------------------------------------------------------------------------------------------------")
		// fmt.Printf("Got block for piece %d of length %d beginning at %d from peer %v\n", payload.Index, len(payload.Block), payload.Begin, pp.client.Peer.IP)
	default:
		
	}

	return nil
}

const maxBlockSize = 16384

func download(client *client.Client, pw *piece.Work) ([]byte, error) {
	buf := make([]byte, pw.Length)
	prog := &pieceProgress{
		buf: buf,
		client: client,
		requested: 0,
		downloaded: 0,
	}

	
	for prog.downloaded < pw.Length {
		if prog.requested < pw.Length && !prog.client.Choked {
			// Determine block size
			blockSize := maxBlockSize

			// If only a part of the max size if left then choose that part as size
			if pw.Length - prog.requested < maxBlockSize {
				blockSize = pw.Length - prog.requested
			}

			// Request for a piece
			err := client.SendRequest(blockSize, prog.requested, pw.Index)
			if err != nil {
				return nil, err
			}
			prog.requested += blockSize
			//fmt.Printf("Requested %d bytes of piece %d from peer %v\n", blockSize, pw.Index, prog.client.Peer.IP)
		}

		msg, err := message.Read(client.Conn)
		if err != nil {
			return nil, err
		}

		prog.handleMessage(msg)
	}
	

	return buf, nil
}

// StartWorker starts a download from the peer that sends results to the process' channel
func StartWorker(client *client.Client, workque chan *piece.Work, results chan *piece.Result) (error) {
	defer client.Conn.Close()

	if err := client.SendUnchoke(); err != nil {
		return err
	}
	// Inform that we are interested in downloading
	if err := client.SendInterested(); err != nil {
		return err
	}


	for pw := range workque {
		// Check if peer has a piece
		haspiece := client.BitField.HasPiece(int(pw.Index))
		if !haspiece {
			fmt.Printf("Peer %v did NOT have the piece of index %d. Pulling a new one...\n", client.Peer.IP, pw.Index)
			workque <- pw
			continue
		}
		//fmt.Printf("Peer %v had the piece of index %d. Starting a download...\n", client.Peer.IP, pw.Index)

		_, err := download(client, pw)
		if err != nil {
			workque <- pw
			return err
		}

		log.Printf("Successfully downloaded piece %d\n", pw.Index)
	}

	return nil
}