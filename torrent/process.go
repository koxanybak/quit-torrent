package torrent

import "math/rand"

// Process represents a torrent process
type Process struct {
	PeerID			[20]byte
	Torrent			TorrentFile
	Peers			[]Peer
	paused			bool
}

// Start starts the torrent process
func (p *Process) Start() {
	
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