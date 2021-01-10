package peerclient

import (
	"net"
)

type msgID uint8
const (
	msgChoke			msgID = 0
	msgUnChoke			msgID = 1
	msgInterested		msgID = 2
	msgNotInterested	msgID = 3
	msgHave				msgID = 4
	msgBitfield			msgID = 5
	msgRequest			msgID = 6
	msgPiece			msgID = 7
	msgCancel			msgID = 8
)

// DownloadWorker represents a download worker. Who would have guessed...
type DownloadWorker struct {
	net.Conn
}

type message struct {
	ID			msgID
	Payload		[]byte
}
