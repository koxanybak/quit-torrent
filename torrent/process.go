package torrent

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/koxanybak/quit-torrent/torrent/client"
	"github.com/koxanybak/quit-torrent/torrent/download"
	"github.com/koxanybak/quit-torrent/torrent/file"
	"github.com/koxanybak/quit-torrent/torrent/peer"
	"github.com/koxanybak/quit-torrent/torrent/piece"
	"github.com/koxanybak/quit-torrent/torrent/saver"
)

// Process represents a torrent process
type Process struct {
	PeerID			[20]byte
	Torrent			file.TorrentFile
	Peers			[]peer.Peer

	paused			bool
}

// Start starts the torrent process
func (p *Process) Start() {
	//done := make(chan struct{})
	workque := p.genPieceWorks()
	results := make(chan *piece.Result)

	for i := range p.Peers {
		peer := p.Peers[i]

		go func() {
			defer func() {
				// if err := recover(); err != nil {
				// 	log.Println("Recovered from panic:", err)
				// }
				log.Println("Number of active downloads: ", runtime.NumGoroutine() - 1)
			}()

			client, err := client.New(p.Torrent.InfoHash, p.PeerID, peer)
			if err != nil {
				log.Printf("Failed to establish connection with peer %v because %v\n", peer.IP, err)
				return
			}

			// TODO: EOF when establishing handshake and keep alive message from peer nil reference

			err = download.StartWorker(client, workque, results)
			if err != nil {
				log.Printf("Downloading failed with peer %v because %v\n", peer.IP, err)
				return
			}
		}()
	}

	for res := range results {
		if err := saver.Save(res); err != nil {
			log.Panicf("Error saving downloaded piece %d\n", res.Index)
		}
		log.Printf("Successfully downloaded piece %d with %d peers\n", res.Index, runtime.NumGoroutine() - 1)
	}
}

// NewProcess creates a new paused torrent process
func NewProcess (filepath string) (*Process, error) {
	torFile, err := file.Open(filepath)
	if err != nil {
		return nil, err
	}

	peerID := [20]byte{}
	// rand.Seed(1)
	_, err = rand.Read(peerID[:])
	if err != nil {
		return nil, err
	}
	
	peers, err := peer.GetPeers(torFile, peerID)
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

func (p *Process) genPieceWorks() chan *piece.Work {
	workque := make(chan *piece.Work, len(p.Torrent.PieceHashes))
	for i, hash := range p.Torrent.PieceHashes {
		piece := piece.NewWork(&p.Torrent, i, hash)
		workque <- piece
	}

	return workque
}