package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Block struct {
	PrevHash   []byte
	Generation uint64
	Difficulty uint8
	Data       string
	Proof      uint64
	Hash       []byte
}

// Create new initial (generation 0) block.
func Initial(difficulty uint8) Block {
	// TODO
	prevHash := []byte(strings.Repeat("\x00", 32))
	gen := uint64(0)
	data := ""
	var p uint64
	var h []byte

	b := Block{prevHash, gen, difficulty, data, p, h}
	return b
}

// Create new block to follow this block, with provided data.
func (prev_block Block) Next(data string) Block {
	// TODO
	prevHash := prev_block.Hash
	gen := prev_block.Generation + 1
	dif := prev_block.Difficulty
	var p uint64
	var h []byte

	b := Block{prevHash, gen, dif, data, p, h}
	return b
}

// Calculate the block's hash.
func (blk Block) CalcHash() []byte {
	// TODO

	prevHashStr := hex.EncodeToString(blk.PrevHash)
	gen := blk.Generation
	dif := blk.Difficulty
	data := blk.Data
	prf := blk.Proof

	hashStr := fmt.Sprintf("%s:%d:%d:%s:%d", prevHashStr, gen, dif, data, prf)
	// fmt.Println(hashStr)

	hash := sha256.New()

	hash.Write([]byte(hashStr))

	return hash.Sum(nil)
}

func (b Block) printHash() {
	fmt.Println(hex.EncodeToString(b.CalcHash()))
}

// Is this block's hash valid?
func (blk Block) ValidHash() bool {
	// TODO

	hash := blk.Hash

	if hash == nil {
		return false
	}

	l := len(hash)

	nBytes := int(blk.Difficulty / 8)
	nBits := blk.Difficulty % 8

	for i := l-1; i >= l-nBytes; i-- {
		if hash[i] != '\x00' {
			return false
		}
	}

	if hash[l-nBytes-1] % (1<<nBits) != 0 {
		return false
	}

	return true
}

// Set the proof-of-work and calculate the block's "true" hash.
func (blk *Block) SetProof(proof uint64) {
	blk.Proof = proof
	blk.Hash = blk.CalcHash()
}

func (b *Block) FindProof() {
	i := uint64(0)
	for ; !(b.ValidHash()); i++ {
		b.SetProof(i)
	}

	fmt.Println(i)

}

func BlockTest() {
	b0 := Initial(16)
	b0.FindProof()

	b0.printHash()
}
