package client

// BitField tells what pieces the peer has
type BitField []byte

// HasPiece tells if the bitfield contains the piece at the index
func (bf BitField) HasPiece(index int) bool {
	// Fancy bit manipulation
	byteIndex := index / 8
	bitIndex := index % 8

	return bf[byteIndex]>>(7-bitIndex)&1 == 1
}

// SetPiece sets the piece defined by the index to true.
// This should be called when getting a Have message from a peer.
func (bf BitField) SetPiece(index int) {
	// Fancy bit manipulation
	byteIndex := index / 8
	bitIndex := index % 8

	bf[byteIndex] = bf[byteIndex] | (1<<(7-bitIndex))
}