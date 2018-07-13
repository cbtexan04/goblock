package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Block struct {
	BPM       int
	Timestamp string
	Hash      string
	PrevHash  string
	Index     int
}

var (
	Blockchain []*Block
	bcMutex    = &sync.Mutex{}
)

func (b *Block) CalculateHash() string {
	// Calculate a unique hash based on the contents of the block
	str := fmt.Sprintf("%d%s%d%s", b.Index, b.Timestamp, b.BPM, b.PrevHash)

	h := sha256.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

// Generates a new block
func (b *Block) Generate(BPM int) (*Block, error) {
	newBlock := Block{
		Index:     b.Index + 1,
		Timestamp: time.Now().String(),
		BPM:       BPM,
		PrevHash:  b.Hash,
	}

	newBlock.Hash = newBlock.CalculateHash()

	return &newBlock, nil
}

func (b *Block) IsValidNextBlock(nextBlock *Block) bool {
	// Make sure the block is incremented correctly
	if b.Index+1 != nextBlock.Index {
		return false
	}

	// Make sure the previous hash of the next block matches our current block
	if b.Hash != nextBlock.PrevHash {
		return false
	}

	// Make sure the hash of the next block is valid
	if nextBlock.CalculateHash() != nextBlock.Hash {
		return false
	}

	return true
}
