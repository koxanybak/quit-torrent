package torrent

import (
	"fmt"
	"net"
	"time"
)


func (p *Process) download(peer Peer) (error) {
	fmt.Println("0")
	conn, err := net.DialTimeout("tcp", peer.String(), 4*time.Second)
	fmt.Println("1")
	if err != nil {
		return err
	}
	defer conn.Close()

	sentHandshake := handshake{
		Pstr: pstr,
		InfoHash: p.Torrent.InfoHash,
		PeerID: p.PeerID,
	}
	fmt.Println("2")
	_, err = conn.Write(sentHandshake.Serialize())
	if err != nil {
		return err
	}
	fmt.Println("3")

	if err := getAndValidateHandshakeFromPeer(conn, &sentHandshake); err != nil {
		return err
	}
	fmt.Println("4")

	return nil
}


type handshake struct {
	Pstr			string
	InfoHash		[20]byte
	PeerID			[20]byte
}
var pstr string = "BitTorrent protocol"

func (h *handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr) + 49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])

	return buf
}

func getAndValidateHandshakeFromPeer(conn net.Conn, sentHandshake *handshake) (error) {
	buf := make([]byte, len(pstr) + 49)
	_, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if len(buf) != len(pstr) + 49 {
		return fmt.Errorf("Wrong length handshake from peer: %d", len(buf))
	}

	hsFromPeer := handshake{}
	hsFromPeer.Pstr = string(buf[1:1+len(pstr)])
	copy(hsFromPeer.InfoHash[:], buf[9+len(pstr) : 29+len(pstr)])
	copy(hsFromPeer.PeerID[:], buf[29+len(pstr) : 49+len(pstr)])

	if hsFromPeer.InfoHash != sentHandshake.InfoHash {
		return fmt.Errorf("Malformed info hash from peer")
	}
	if hsFromPeer.Pstr != sentHandshake.Pstr {
		return fmt.Errorf("Malformed protocol from peer: %s", hsFromPeer.Pstr)
	}

	return nil
}