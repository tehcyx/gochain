package blockchain

import (
	"fmt"
	"testing"
)

// Tests
func TestNewBlockchain(t *testing.T) {
	chain := NewBlockchain()
	if len(chain.Chain) != 1 {
		t.Errorf("Blockchain should be initialized with Genesis block upon creation and therefor length should be 1")
		t.Fail()
	}
	if len(chain.CurrentTransactions) > 0 {
		t.Errorf("Transactions currently queued in a newly generated blockchain should be 0")
		t.Fail()
	}
	if len(chain.Nodes) > 0 {
		t.Errorf("Nodes currently registered to a newly generated blockchain should be 0")
		t.Fail()
	}
}

func TestProofOfWork(t *testing.T) {
	chain := NewBlockchain()
	number := chain.ProofOfWork(100)
	if number != 29031 {
		t.Errorf("Something changed with the difficulty. Check that it is still \"0000\"")
		t.Fail()
	}
}

func TestStringBlock(t *testing.T) {
	chain := NewBlockchain()
	block := chain.LastBlock()
	index := block.Index
	hash := block.Hash
	previousHash := block.PreviousHash
	proof := block.Proof
	timestamp := block.Timestamp
	transactions := block.Transactions
	expected := fmt.Sprintf("Index: %d,\n Hash: %s,\n PreviousHash: %s,\n Proof: %d,\n Timestamp: %d,\n Transactions: %v", index, hash, previousHash, proof, timestamp, transactions)
	if block.String() != expected {
		t.Errorf("ToString should return a different result")
		t.Fail()
	}
}

func TestStringBlockchain(t *testing.T) {
	blockchain := NewBlockchain()

	chain := blockchain.Chain
	currentTransactions := blockchain.CurrentTransactions
	nodes := blockchain.Nodes

	expected := fmt.Sprintf("Chain: % +v,\n Current Transactions: % +v,\n Nodes: % +v", chain, currentTransactions, nodes)
	if blockchain.String() != expected {
		t.Errorf("ToString should return a different result")
		t.Fail()
	}
}

func TestNewBlock(t *testing.T) {
	chain := NewBlockchain()
	originalBlock := chain.LastBlock()
	block := chain.NewBlock(29031, originalBlock.Hash)

	if block.Hash == "" {
		t.Errorf("A block must not come without a hash")
		t.Fail()
	}
	if block.PreviousHash != originalBlock.Hash {
		t.Errorf("Previous hash has to be set to original blocks hash")
		t.Fail()
	}
	if block.Index == originalBlock.Index {
		t.Errorf("the latest created block should not have the same index")
		t.Fail()
	}
	if block.Index != originalBlock.Index+1 {
		t.Errorf("the latest created block should have the index of the original block + 1")
		t.Fail()
	}
	if len(block.Transactions) > 0 {
		t.Errorf("No transactions should be included in this one")
		t.Fail()
	}
	chain = NewBlockchain()
	blockNumber := chain.NewTransaction("me", "you", 1.2)
	block = chain.NewBlock(29031, originalBlock.Hash)
	if block.Index != blockNumber {
		t.Errorf("newly mined block should be the next one, that the transaction will be included")
		t.Fail()
	}
	if len(block.Transactions) != 1 {
		t.Errorf("Transactions in this block should be exactly one")
		t.Fail()
	}
}

func TestHash(t *testing.T) {
	chain := NewBlockchain()
	block := chain.NewBlock(29031, chain.Chain[0].PreviousHash)

	if chain.Hash(*block) != block.Hash {
		t.Errorf("A blocks hash should always be the same")
		t.Fail()
	}
}

func TestLastBlock(t *testing.T) {
	chain := NewBlockchain()
	firstBlock := chain.Chain[0]
	lastBlock := chain.LastBlock()
	if firstBlock.Hash != lastBlock.Hash {
		t.Errorf("First Block is the only block, so they should be equal")
		t.Fail()
	}
	chain.NewBlock(29031, lastBlock.Hash)
	lastBlock = chain.LastBlock()
	if firstBlock.Hash == lastBlock.Hash {
		t.Errorf("First Block is not the only block anymore, so they must not be equal")
		t.Fail()
	}
}

func TestRegisterNode(t *testing.T) {
	chain := NewBlockchain()
	isRegistered := chain.RegisterNode("http://localhost:8080", "just a node")
	if !isRegistered {
		t.Errorf("http://localhost:8080 should be a valid url to be registered")
		t.Fail()
	}
	isRegistered = chain.RegisterNode("http://google.com", "another node")
	if !isRegistered {
		t.Errorf("http://google.com should be a valid url to be registered")
		t.Fail()
	}
	if len(chain.Nodes) != 2 {
		t.Errorf("List of registered nodes should be exactly 2")
		t.Fail()
	}
	isRegistered = chain.RegisterNode("not a url", "not working")
	if isRegistered {
		t.Errorf("\"not a url\" should not register as node")
		t.Fail()
	}
}

func TestValidProof(t *testing.T) {
	chain := NewBlockchain()
	isValid := chain.ValidProof(0, 0)
	if isValid {
		t.Errorf("ValidProof of 0,0 should be false")
		t.Fail()
	}
	isValid = chain.ValidProof(100, 29031)
	if !isValid {
		t.Errorf("ValidProof of 100, 29031 should be true")
		t.Fail()
	}
}

// Benchmarks
func BenchmarkProofOfWork(b *testing.B) {
	chain := NewBlockchain()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain.ProofOfWork(chain.LastBlock().Proof)
	}
}
