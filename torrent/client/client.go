package client

import (
	"fmt"
	"net"
	"time"

	"github.com/koxanybak/quit-torrent/torrent/message"
	"github.com/koxanybak/quit-torrent/torrent/peer"
)

// Client represents a client that can talk to a peer
type Client struct {
	Conn		net.Conn
	Peer		peer.Peer
	Choked		bool
	BitField	BitField
}

// SendUnchoke s
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
// SendInterested s
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
// SendRequest s
func (c *Client) SendRequest(blockSize int, begin int, index int) error {
	msg := message.Message{
		ID: message.MsgRequest,
		Payload: message.CreateRequestPayload(blockSize, begin, index),
	}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// New opens a BitTorrent client to a peer (does handshake and gathers the bitfield)
// Remember to close the connection when done
func New(infoHash [20]byte, peerID [20]byte, peer peer.Peer) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 4*time.Second)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Conn: conn,
		Peer: peer,
		Choked: true,
	}

	// Handshake
	if err := client.DoHandshake(infoHash, peerID); err != nil {
		conn.Close()
		return nil, fmt.Errorf("Error completing handshake: %v", err)
	}

	// Get bitfield
	msg, err := message.Read(client.Conn)
	if err != nil || msg.ID != message.MsgBitfield {
		conn.Close()
		return nil, fmt.Errorf("Peer didn't send bitfield: %v", err)
	}
	client.BitField = msg.Payload

	return client, err
}
