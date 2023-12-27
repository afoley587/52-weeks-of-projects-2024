package main

import (
	"log"
	"time"
)

/*
 * Our Blockchain Struct will have a few attributes:
 *		1. The genesis block which starts this chain
 *		2. The actual chain itself
 *		3. The difficulty of the chain (for mining more)
 */
type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

/*
 * When we create a new blockchain, we will need to
 * first create its genesis block. We will then box up
 * the genesis block, new chain, and chain difficulty
 * into a struct which our peers can pass around later.
 */
func NewBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		Hash:      "0",
		Height:    0,
		Timestamp: time.Now().Unix(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

/*
 * We also need a facility to add some structured data to our blockchain.
 * Our appendBlock function will do just that. It will first take a location
 * and a wave height from the caller. It will then create a block with that data,
 * assign it a height, and calculate a hash for it. Then, it will add it to the
 * blockchain.
 */
func (b *Blockchain) appendBlock(location string, waveHeight int) {
	blockData := BlockData{
		Location:   location,
		WaveHeight: waveHeight,
	}
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         blockData,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().Unix(),
		Height:       lastBlock.Height + 1,
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

/*
 * The validity of chains is one of the main reason to use them. They're
 * secure because it should be impossible to tamper with them. So, how do
 * we tell if we are maintaining our validity and integrity? Well, we will
 * travers our blockchain. For each block, we can do a few checks:
 * 		1. Do the heights make sense? Is the previous blocks height 1 less
 *			than the current block's height.
 *		2. Do the hashes make sense? Is the current block's the same as it
 *			originally was when i calculated it?
 *		3. Do the linkages make sense? Is my current block's previous hash
 *			actually the same as the previous block's hash?
 * If the answer to any of these is no, then we have an issue and our chain
 * has become invalid somewhere.
 */
func (b Blockchain) isValid() bool {
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		if currentBlock.Height != previousBlock.Height+1 {
			log.Println("Bad Height")
			return false
		}
		if currentBlock.Hash != currentBlock.calculateHash() {
			log.Println("Bad Hash")
			return false
		}
		if currentBlock.PreviousHash != previousBlock.Hash {
			log.Println("Bad Prev Hash")
			return false
		}
	}
	return true
}
