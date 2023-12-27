package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

/*
 * Our blocks will hold surfing data. We will be keeping
 * track of the recorded locations and the recorded wave heights
 */
type BlockData struct {
	Location   string
	WaveHeight int
}

/*
 * Our Block Struct will have a few attributes:
 *		1. Data to record on the blockchain
 *		2. A block hash, the ID of the block generated using cryptography
 *			techniques
 *		3. The previous block’s hash is the cryptographic hash of the
 *			last block in the blockchain. It is recorded in every block to
 *			link it to the chain and improve its security
 *		4. A timestamp of when the block was created and added to the
 *			blockchain
 *		5. Height is the index of the block on the chain
 *		6. Proof Of Work shows how much the miner had to work to achieve
 *			an acceptable block
 */
type Block struct {
	Data         BlockData
	Hash         string
	PreviousHash string
	Timestamp    int64
	Height       int
	Pow          int
}

/*
 * A hash of a block is it's ID. It should be 100% unique across the entire
 * blockchain. We can compute this in a number of ways (as long as it's unique),
 * but we will combine and hash the following pieces of data:
 *		1. the previous block's hash
 *		2. the stringified data of the block
 *		3. the current timestamp
 *		4. the block height
 */
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + strconv.FormatInt(b.Timestamp, 10) + strconv.Itoa(b.Height) + strconv.Itoa(b.Pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

/*
 * Mining is essentially adding new blocks to our block chain with a certain
 * difficulty. In the context of blockchain, difficulty refers to a parameter
 * that regulates how challenging it is to add a new block to the blockchain.
 * The difficulty level is dynamically adjusted to ensure that the average time
 * between the creation of new blocks remains relatively constant. This is crucial for
 * maintaining the consistency and security of the blockchain.
 * The difficulty is usually set in such a way that miners,
 * who are participants in the network responsible for validating
 * and adding new blocks, need to solve a complex mathematical problem to
 * create a new block. The difficulty adjusts regularly based on factors such as
 * the total computational power of the network. If more miners join the
 * network and the overall computational power increases, the difficulty level
 * is raised to maintain a consistent block creation time.
 */
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Pow++
		b.Hash = b.calculateHash()
	}
}
