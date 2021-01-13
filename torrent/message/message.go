package message

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// MsgID is the id of BitTorrent messages
type MsgID uint8
const (
	MsgChoke			MsgID = 0
	MsgUnchoke			MsgID = 1
	MsgInterested		MsgID = 2
	MsgNotInterested	MsgID = 3
	MsgHave				MsgID = 4
	MsgBitfield			MsgID = 5
	MsgRequest			MsgID = 6
	MsgPiece			MsgID = 7
	MsgCancel			MsgID = 8
	MsgKeepAlive		MsgID = 9
)

func (mid MsgID) String() string {
	if mid == MsgChoke {
		return "Choke"
	}
	if mid == MsgUnchoke {
		return "Unchoke"
	}
	if mid == MsgInterested {
		return "Interested"
	}
	if mid == MsgNotInterested {
		return "NotInterested"
	}
	if mid == MsgHave {
		return "Have"
	}
	if mid == MsgBitfield {
		return "Bitfield"
	}
	if mid == MsgRequest {
		return "Request"
	}
	if mid == MsgPiece {
		return "Piece"
	}
	if mid == MsgCancel {
		return "Cancel"
	}

	if mid == MsgKeepAlive {
		return "KeepAlive"
	}

	return fmt.Sprintf("UNKNOWN_ID_%d", mid)
}

// Message represents a message sent/received over the BitTorrent connection
type Message struct {
	ID			MsgID
	Payload		[]byte
}

// Read reads a message from the peer
func Read(conn net.Conn) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		msg := &Message{
			ID: MsgKeepAlive,
			Payload: nil,
		}
		return msg, nil // Keep alive
	}
	
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(conn, messageBuf)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		ID: MsgID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return msg, nil
}

// Serialize turns message into []byte
func (m *Message) Serialize() []byte {
	length := len(m.Payload) + 1
	buf := make([]byte, length + 4)

	binary.BigEndian.PutUint32(buf[:4], uint32(length))
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)

	return buf
}

// CreateRequestPayload creates a payload for a 'request' message
func CreateRequestPayload(blockSize int, begin int, index int) []byte {
	buf := make([]byte, 12)

	binary.BigEndian.PutUint32(buf[0:4], uint32(index))
	binary.BigEndian.PutUint32(buf[4:8], uint32(begin))
	binary.BigEndian.PutUint32(buf[8:12], uint32(blockSize))

	return buf
}
