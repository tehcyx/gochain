package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

// Blockchain struct
type Blockchain struct {
	Chain               []Block           `json:"chain"`
	CurrentTransactions []Transaction     `json:""`
	Nodes               map[string]string `json:"nodes"`
}

// Block struct
type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int64         `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
	Hash         string        `json:"hash"`
}

// Transaction struct
type Transaction struct {
	Sender    string  `json:"sender"`
	Recepient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

// ProofOfWork create a proof of work
func (c *Blockchain) ProofOfWork(lastProof int64) int64 {
	proof := int64(0)
	for !c.ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

// Hash create a hash of a block
func (c *Blockchain) Hash(b Block) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d %s %d %d %v", b.Index, b.PreviousHash, b.Proof, b.Timestamp, b.Transactions)))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	return sha256Hash
}

// NewBlock create a new block
func (c *Blockchain) NewBlock(proof int64, previousHash string) *Block {
	b := new(Block)
	b.Index = len(c.Chain) + 1
	b.Timestamp = time.Now().Unix()
	b.Proof = proof
	b.PreviousHash = previousHash
	b.Transactions = c.CurrentTransactions
	b.Hash = c.Hash(*b)

	// reset current transactions
	c.CurrentTransactions = []Transaction{}

	c.Chain = append(c.Chain, *b)
	return b
}

// NewTransaction creates a new transaction that will be part of the next mined block
func (c *Blockchain) NewTransaction(sender, recipient string, amount float64) int {
	t := new(Transaction)
	t.Sender = sender
	t.Recepient = recipient
	t.Amount = amount

	c.CurrentTransactions = append(c.CurrentTransactions, *t)

	return c.LastBlock().Index + 1
}

// LastBlock returns last block of chain
func (c *Blockchain) LastBlock() Block {
	return c.Chain[len(c.Chain)-1]
}

// NewBlockchain init the blockchain
func NewBlockchain() *Blockchain {
	blockchain := new(Blockchain)
	blockchain.CurrentTransactions = []Transaction{}
	blockchain.Nodes = map[string]string{}

	blockchain.NewBlock(100, "1")
	return blockchain
}

// RegisterNode registers given node with given comment
func (c *Blockchain) RegisterNode(address, comment string) bool {
	parsedURL, err := url.ParseRequestURI(address)
	if err != nil {
		return false
	}
	c.Nodes[parsedURL.Host] = comment
	return true
}

// ValidProof checks if the proof given is valid
func (c *Blockchain) ValidProof(lastProof int64, proof int64) bool {
	h := sha256.New()
	h.Write([]byte(string(lastProof) + string(proof)))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	return sha256Hash[:4] == "0000"
}

func (c Blockchain) String() string {
	return fmt.Sprintf("Chain: % +v,\n Current Transactions: % +v", c.Chain, c.CurrentTransactions)
}

func (b Block) String() string {
	return fmt.Sprintf("Index: %d,\n Hash: %s,\n PreviousHash: %s,\n Proof: %d,\n Timestamp: %d,\n Transactions: %v", b.Index, b.Hash, b.PreviousHash, b.Proof, b.Timestamp, b.Transactions)
}

func (t Transaction) String() string {
	return fmt.Sprintf("Sender: %s,\n Recipient: %s,\n Amount: %f", t.Sender, t.Recepient, t.Amount)
}
