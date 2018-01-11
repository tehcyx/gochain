package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Blockchain struct
type Blockchain struct {
	chain               []Block
	currentTransactions []Transaction
}

// Block struct
type Block struct {
	index        int
	timestamp    int64
	transactions []Transaction
	proof        int64
	previousHash string
	hash         string
}

// Transaction struct
type Transaction struct {
	sender    string
	recepient string
	amount    float64
}

// ProofOfWork create a proof of work
func (c *Blockchain) ProofOfWork(lastProof int64) int64 {
	proof := int64(0)
	for !c.validProof(lastProof, proof) {
		proof++
	}
	return proof
}

// Hash create a hash of a block
func (c *Blockchain) Hash(b Block) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d %s %d %d %v", b.index, b.previousHash, b.proof, b.timestamp, b.transactions)))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	return sha256Hash
}

// NewBlock create a new block
func (c *Blockchain) NewBlock(proof int64, previousHash string) *Block {
	b := new(Block)
	b.index = len(c.chain) + 1
	b.timestamp = time.Now().Unix()
	b.proof = proof
	b.previousHash = previousHash
	b.transactions = c.currentTransactions
	b.hash = c.Hash(*b)

	// reset current transactions
	c.currentTransactions = []Transaction{}

	c.chain = append(c.chain, *b)
	return b
}

// NewTransaction creates a new transaction that will be part of the next mined block
func (c *Blockchain) NewTransaction(sender, recipient string, amount float64) int {
	t := new(Transaction)
	t.sender = sender
	t.recepient = recipient
	t.amount = amount

	c.currentTransactions = append(c.currentTransactions, *t)

	return c.LastBlock().index + 1
}

// LastBlock returns last block of chain
func (c *Blockchain) LastBlock() Block {
	return c.chain[len(c.chain)-1]
}

func (b *Block) GetProof() int64 {
	return b.proof
}

func (b *Block) GetPreviousHash() string {
	return b.previousHash
}

// NewBlockchain init the blockchain
func NewBlockchain() *Blockchain {
	blockchain := new(Blockchain)
	blockchain.currentTransactions = []Transaction{}

	blockchain.NewBlock(100, "1")
	return blockchain
}

func (c *Blockchain) validProof(lastProof int64, proof int64) bool {
	h := sha256.New()
	h.Write([]byte(string(lastProof) + string(proof)))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	return sha256Hash[:4] == "0000"
}

func (c Blockchain) String() string {
	return fmt.Sprintf("Chain: % +v,\n Current Transactions: % +v", c.chain, c.currentTransactions)
}

func (b Block) String() string {
	return fmt.Sprintf("Index: %d,\n Hash: %s,\n PreviousHash: %s,\n Proof: %d,\n Timestamp: %d,\n Transactions: %v", b.index, b.hash, b.previousHash, b.proof, b.timestamp, b.transactions)
}

func (t Transaction) String() string {
	return fmt.Sprintf("Sender: %s,\n Recipient: %s,\n Amount: %f", t.sender, t.recepient, t.amount)
}
