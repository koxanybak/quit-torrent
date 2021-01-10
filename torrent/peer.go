package torrent

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"

	"github.com/zeebo/bencode"
)

type todo struct {
	Complete		uint64		`bencode:"complete"`
	Downloaded		uint64		`bencode:"downloaded"`
	Incomplete		uint64		`bencode:"incomplete"`
	Interval		uint64		`bencode:"interval"`
	MinInterval		uint64		`bencode:"min interval"`
	Peers			[]byte		`bencode:"peers"`
}

// Peer represents a peer in the torrent network
type Peer struct {
	IP		net.IP
	Port	uint16
}

func (p *Peer) String() string {
	return fmt.Sprintf("%s:%d", p.String(), p.Port)
}

func unmarshalPeers(peersBin []byte) ([]Peer, error) {
	peerSize := 6
	if len(peersBin) % peerSize != 0 {
		return nil, fmt.Errorf("Malformed peer list. Size: %d", len(peersBin))
	}

	numPeers := len(peersBin) / peerSize
	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		peers[i] = Peer{
			IP: peersBin[i*peerSize : i*peerSize + 4],
			Port: binary.BigEndian.Uint16(peersBin[i*peerSize + 4 : i*peerSize + 6]),
		}
	}

	return peers, nil
}

// GetPeers returns
func GetPeers(t *TorrentFile, peerID [20]byte) ([]Peer, error) {
	url, err := t.GetTrackerURL(peerID, 8001)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trackerData todo
	if err := bencode.NewDecoder(resp.Body).Decode(&trackerData); err != nil {
		return nil, err
	}

	peers, err := unmarshalPeers(trackerData.Peers)
	if err != nil {
		return nil, err
	}

	return peers, nil
}