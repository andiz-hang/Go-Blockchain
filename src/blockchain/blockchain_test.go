package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"bytes"
	"strings"
	_ "encoding/hex"
	_ "fmt"
)

// TODO: some useful tests of Blocks

func TestInitial(t *testing.T) {
	b := Initial(16)

	assert.True(t, bytes.Equal(b.PrevHash, []byte(strings.Repeat("\x00", 32))), "Initial block's prevHash incorrect")
	assert.Equal(t, b.Generation, uint64(0), "Initial block's generation is not 0")
	assert.Equal(t, b.Difficulty, uint8(16), "Initial block's difficulty is not set properly")
	assert.Equal(t, b.Data, "", "Initial Block's data not empty string")
}

func TestNext(t *testing.T) {
	b := Initial(16)
	next := b.Next("Some Data")
	
	assert.True(t, bytes.Equal(b.Hash, next.PrevHash), "Next prevHash incorrect")
	assert.Equal(t, b.Generation + 1, next.Generation, "Next Generation incorrect")
	assert.Equal(t, b.Difficulty, next.Difficulty, "Next Difficulty incorrect")
	assert.Equal(t, "Some Data", next.Data, "Next sets Data incorrectly")
}

func TestValidHash(t *testing.T) {
	b := Initial(16)
	assert.False(t, b.ValidHash())

	b.Hash = []byte(strings.Repeat("\x00", 32))
	assert.True(t, b.ValidHash())

	b.Hash = []byte("\x01")
	assert.False(t, b.ValidHash())

	b = Initial(19)
	b.Hash = []byte("\x08\x00\x00")
	assert.True(t, b.ValidHash())

	b.Hash = []byte("\x07\x00\x00")
	assert.False(t, b.ValidHash())
}

func TestMineRange(t *testing.T) {
	b := Initial(16)

	// Should return no solution
	res := b.MineRange(0, 50000, 20, 2000)
	assert.Equal(t, false, res.Found)

	res = b.MineRange(0, 56231, 20, 2000)
	assert.Equal(t, true, res.Found)
	assert.Equal(t, uint64(56231), res.Proof)

	res = b.MineRange(56231, 60000, 10, 100)
	assert.Equal(t, true, res.Found)
	assert.Equal(t, uint64(56231), res.Proof)
}

func TestIsValid(t *testing.T) {
	c := Blockchain{}
	init := Initial(16)
	init.Mine(20)

	c.Add(init)
	assert.True(t, c.IsValid())

	next := init.Next("Spooky Message")
	next.Mine(20)

	c.Add(next)
	assert.True(t, c.IsValid())

	next = next.Next("Bitcoin is the future")
	next.Mine(20)

	c.Add(next)
	assert.True(t, c.IsValid())
}