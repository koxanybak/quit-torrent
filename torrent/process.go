package torrent

import (
	"log"
	"math/rand"
)

// Process represents a torrent process
type Process struct {
	PeerID			[20]byte
	Torrent			TorrentFile
	Peers			[]Peer

	paused			bool

	done			chan struct{}
}

// Start starts the torrent process
func (p *Process) Start() {
	// done := make(chan struct{})
	// data := make(chan []byte)
	log.Println("Juu")
	p.download(p.Peers[0])
	log.Println("Moi")
}

// NewProcess creates a new paused torrent process
func NewProcess (filepath string) (*Process, error) {
	torFile, err := Open(filepath)
	if err != nil {
		return nil, err
	}

	peerID := [20]byte{}
	rand.Seed(1)
	rand.Read(peerID[:])
	
	peers, err := GetPeers(torFile, peerID)
	if err != nil {
		return nil, err
	}

	return &Process{
		PeerID: peerID,
		Torrent: *torFile,
		Peers: peers,
	}, nil
}