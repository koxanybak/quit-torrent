package torrent

import (
	"math/rand"
	"net"
	"time"

	"github.com/koxanybak/quit-torrent/peerclient"
)

// Process represents a torrent process
type Process struct {
	PeerID			[20]byte
	Torrent			TorrentFile
	Peers			[]Peer

	paused			bool

	done			chan struct{}
}

// startDownloadWorker starts a download from the peer that sends results to the process' channel
func (p *Process) startDownloadWorker(peer Peer) (error) {
	var worker peerclient.DownloadWorker
	var err error
	worker.Conn, err = net.DialTimeout("tcp", peer.String(), 4*time.Second)
	if err != nil {
		return err
	}
	defer worker.Close()

	if err := worker.DoHandshake(p.Torrent.InfoHash, p.PeerID); err != nil {
		return err
	}

	return nil
}

// Start starts the torrent process
func (p *Process) Start() {
	// done := make(chan struct{})
	// data := make(chan []byte)
	err := p.startDownloadWorker(p.Peers[0])
	if err != nil {
		panic(err)
	}
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
		paused: true,
	}, nil
}