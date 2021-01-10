package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/zeebo/bencode"
)

type benTorrentInfo struct {
	Length			int64					`bencode:"length"`
	Name			string					`bencode:"name"`
	PieceLength		int64					`bencode:"piece length"`
	Pieces			[]byte					`bencode:"pieces"`
}

type benTorrent struct {
	Announce		string					`bencode:"announce"`
	AnnounceList	[]interface{}			`bencode:"announce-list"`
	Comment			string					`bencode:"comment"`
	CreatedBy		string					`bencode:"created by"`
	CreationDate	int64					`bencode:"creation date"`
	Info			benTorrentInfo			`bencode:"info"`
}

func (info *benTorrentInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.NewEncoder(&buf).Encode(&info); err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (info *benTorrentInfo) pieceHashes() ([][20]byte, error) {
	hashLength := 20
	if len(info.Pieces) % hashLength != 0 {
		err := fmt.Errorf("Malformed number of piece bytes: %d", len(info.Pieces))
		return [][20]byte{}, err
	}

	numHashes := len(info.Pieces) / hashLength
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], info.Pieces[i * hashLength : i * hashLength + hashLength])
	}

	return hashes, nil
}

func (b *benTorrent) toTorrentFile() (*TorrentFile, error) {
	hash, err := b.Info.hash()
	if err != nil {
		return nil, err
	}
	pieceHashes, err := b.Info.pieceHashes()
	if err != nil {
		return nil, err
	}

	return &TorrentFile{
		InfoHash: hash,
		PieceHashes: pieceHashes,
		Announce: b.Announce,
		PieceLength: b.Info.PieceLength,
		Name: b.Info.Name,
		Length: b.Info.Length,
	}, nil
}

// TorrentFile contains all the relevant info of the file
type TorrentFile struct {
	Announce	string
	InfoHash	[20]byte
	PieceHashes	[][20]byte
	PieceLength	int64
	Length		int64
	Name		string
}

// GetTrackerURL returns the url of the tracking server
func (t *TorrentFile) GetTrackerURL(peerID [20]byte, port int) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	fmt.Println(len(peerID))
	params := url.Values{
		"info_hash":	[]string{string(t.InfoHash[:])},
		"peer_id":		[]string{string(peerID[:])},
		"port":			[]string{strconv.Itoa(port)},
		"uploaded":		[]string{"0"},
		"downloaded":	[]string{"0"},
		"compact":		[]string{"1"},
		"left":			[]string{strconv.Itoa(int(t.Length))},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

// Open creates a new torrent file object from a file path
func Open(path string) (*TorrentFile, error) {
	// Open
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode
	decoder := bencode.NewDecoder(file)
	var bt benTorrent
	if err := decoder.Decode(&bt); err != nil {
		return nil, err
	}

	torrentFile, err := bt.toTorrentFile()
	if err != nil {
		return nil, err
	}

	return torrentFile, nil
}