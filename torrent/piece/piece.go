package piece

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	"github.com/koxanybak/quit-torrent/torrent/file"
)

// Work represents a piece to be downloaded
type Work struct {
	Length		int
	Hash		[20]byte
	Begin		int
	Index		int
	End			int
}

// NewWork creates a new piece work :DDD
func NewWork(tf *file.TorrentFile, index int, hash [20]byte) *Work {
	begin, end := tf.BoundsForPieceAt(index)

	length := end - begin

	work := &Work{
		Hash: hash,
		Length: length,
		Index: index,
		Begin: begin,
		End: end,
	}

	return work
}

// Payload represents a payload from a 'piece' message
type Payload struct {
	Begin		int
	Index		int
	Block		[]byte
}

// Result represents a downloaded piece
type Result struct {
	Length		int
	Index		int
	Data		[]byte
	Begin		int
	End			int
}

// NewResult creates a new downloaded piece from the 'piecework' and downloaded data
func NewResult(pw *Work, data []byte) *Result {
	return &Result{
		Data: data,
		Begin: pw.Begin,
		End: pw.End,
		Index: pw.Index,
	}
}

// HashesMatch checks if the hashes of the work and result match
func HashesMatch(pw *Work, pr *Result) bool {
	resultHash := sha1.Sum(pr.Data)
	return bytes.Equal(resultHash[:], pw.Hash[:])
}

// UnmarshalPiecePayload turns payload data into a piece
func UnmarshalPiecePayload(payload []byte) (*Payload, error) {
	if len(payload) <= 8 {
		return nil, fmt.Errorf("Too short payload for piece message: %d", len(payload))
	}

	index := binary.BigEndian.Uint32(payload[:4])
	begin := binary.BigEndian.Uint32(payload[4:8])
	block := payload[8:]

	return &Payload{
		Index: int(index),
		Begin: int(begin),
		Block: block,
	}, nil
}
