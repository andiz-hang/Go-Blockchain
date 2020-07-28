package blockchain

import (
	"work_queue"
	"fmt"
	"encoding/hex"
)

type miningWorker struct {
	// TODO. Should implement work_queue.Worker
	b Block
	start, end uint64
}
func (m miningWorker) Run() interface{} {
	res := MiningResult{0, false}

	for i := m.start; i <= m.end; i++ {
		m.b.SetProof(i)
		if m.b.ValidHash() {
			res.Proof = i
			res.Found = true
			break
		}
	}

	return res
}

type MiningResult struct {
	Proof uint64 // proof-of-work value, if found.
	Found bool   // true if valid proof-of-work was found.
}

// Mine the range of proof values, by breaking up into chunks and checking
// "workers" chunks concurrently in a work queue. Should return shortly after a result
// is found.
func (blk Block) MineRange(start uint64, end uint64, workers uint64, chunks uint64) MiningResult {
	// TODO
	miningRange := end - start + 1

	chunkSz := miningRange / chunks
	remainder := miningRange % chunks

	queue := work_queue.Create(uint(workers), uint(chunks))

	rStart := start

	// Create 'chunks' amount of 'workers', each of 'chunkSz' size
	for i := uint64(0); i < chunks - remainder; i++ {
		worker := miningWorker{blk, rStart, rStart + chunkSz - 1}
		queue.Enqueue(worker)

		rStart += chunkSz
	}

	// If there is a remainder, create chunks of 'chunkSz' + 1
	for i := uint64(0); i < remainder; i++ {
		worker := miningWorker{blk, rStart, rStart + chunkSz}
		queue.Enqueue(worker)

		rStart += chunkSz + 1
	}

	res := MiningResult{}
	// Receive the results
	for i := uint64(0); i < chunks; i++ {
		tmp := <- queue.Results
		t2, ok := tmp.(MiningResult)
		if ok {
			res = t2
		}

		if res.Found {
			queue.Shutdown()
			break
		}
	}

	// Proof either found or not found
	return res
}

// Call .MineRange with some reasonable values that will probably find a result.
// Good enough for testing at least. Updates the block's .Proof and .Hash if successful.
func (blk *Block) Mine(workers uint64) bool {
	reasonableRangeEnd := uint64(4 * 1 << blk.Difficulty) // 4 * 2^(bits that must be zero)
	mr := blk.MineRange(0, reasonableRangeEnd, workers, 4321)
	if mr.Found {
		blk.SetProof(mr.Proof)
	}
	return mr.Found
}

func MineTest() {
	b0 := Initial(20)
	b0.Mine(1)
	fmt.Println(b0.Proof, hex.EncodeToString(b0.Hash))
	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	fmt.Println(b1.Proof, hex.EncodeToString(b1.Hash))
	b2 := b1.Next("this is not interesting")
	b2.Mine(1)
	fmt.Println(b2.Proof, hex.EncodeToString(b2.Hash))
}