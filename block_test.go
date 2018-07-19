package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateGivesNextIndex(t *testing.T) {
	genesis := &Block{
		Index:     0,
		Timestamp: time.Now().String(),
		BPM:       0,
		PrevHash:  "",
	}
	genesis.Hash = genesis.CalculateHash()

	// Check #1
	newBlock, err := genesis.Generate(10)
	assertEqual(t, nil, err)
	assertEqual(t, genesis.Index+1, newBlock.Index)

	// Check again to be sure
	anotherBlock, err := newBlock.Generate(20)
	assertEqual(t, nil, err)
	assertEqual(t, newBlock.Index+1, anotherBlock.Index)
}

func TestGenerateSetsPrevHash(t *testing.T) {
	genesis := &Block{
		Index:     0,
		Timestamp: time.Now().String(),
		BPM:       0,
		PrevHash:  "",
	}
	genesis.Hash = genesis.CalculateHash()

	// Check #1
	newBlock, err := genesis.Generate(20)
	assertEqual(t, nil, err)
	assertEqual(t, genesis.Hash, newBlock.PrevHash)

	// Check again to be sure
	anotherBlock, err := newBlock.Generate(30)
	assertEqual(t, nil, err)
	assertEqual(t, newBlock.Hash, anotherBlock.PrevHash)
}

/* expected and actual must be of the same type */
func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected == actual {
		return
	}

	err := fmt.Sprintf("Expected %v but received %v", expected, actual)
	t.Fatal(err)
}
