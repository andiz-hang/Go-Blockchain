package blockchain

import (
	_ "crypto/sha256"
	_ "encoding/hex"
	"fmt"
	"strings"
	"bytes"
)

type Blockchain struct {
	Chain []Block
}

func (chain *Blockchain) Add(blk Block) {
	// TODO
	chain.Chain = append(chain.Chain, blk)
}

func (chain Blockchain) IsValid() bool {
	// TODO
	ch := chain.Chain
	fst := ch[0]

	// The initial block has previous hash all null bytes and is generation zero.
	if !bytes.Equal(fst.PrevHash, []byte(strings.Repeat("\x00", 32))) {
		fmt.Println("Initial Block's PrevHash is not all null bytes")
		return false
	}

	if fst.Generation != 0 {
		fmt.Println("Initial Block's Generation is not 0")
		return false
	}

	dif := fst.Difficulty

	for i := 0; i < len(ch); i++ {
		el := ch[i]

		// Each block has the same difficulty value.
		if el.Difficulty != dif {
			return false
			fmt.Println("Not all blocks have same difficulty")
		}

		if i != 0 {
			prevEl := ch[i-1]
			// Each block has a generation value that is one more than the previous block.
			if el.Generation != prevEl.Generation + 1 {
				fmt.Println("A block has a generation that is not the previous block's generation + 1")
				return false
			}

			// Each block's previous hash matches the previous block's hash.
			if !(bytes.Equal(el.PrevHash, prevEl.Hash)) {
				fmt.Println("A block has a PrevHash that is different from the previous Hash")
				return false
			}
		}

		// Each block's hash value ends in difficulty null bits.
		if !el.ValidHash() {
			fmt.Println("A block doesn't have the correct number of null bits")
			return false
		}

		// Each block's hash value actually matches its contents.

		if !(bytes.Equal(el.CalcHash(), el.Hash)) {
			fmt.Println("A block's hash doesn't match its contents")
			return false
		}
	}

	return true
}
