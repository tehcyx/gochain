package main

import (
	"fmt"

	"github.com/tehcyx/gochain/blockchain"
)

func main() {
	b := blockchain.NewBlockchain()

	i := 0
	for i < 3 {
		i++
		// transactionId := b.NewTransaction("Alice", "Bob", 44.3)
		// fmt.Println(fmt.Sprintf("Transaction will be added to Block %d", transactionId))
		b.NewTransaction("Alice", "Bob", 44.3)
		// transactionId = b.NewTransaction("Bob", "Alice", 12.6)
		// fmt.Println(fmt.Sprintf("Transaction will be added to Block %d", transactionId))
		b.NewTransaction("Bob", "Alice", 12.6)

		lastBlock := b.LastBlock()
		lastProof := lastBlock.GetProof()
		proof := b.ProofOfWork(lastProof)

		// fmt.Println(fmt.Sprintf("Transaction for mining reward of Block %d", b.NewTransaction("0", "sqrt(All Evil)", 1.0)))

		previousHash := b.Hash(lastBlock)
		// block := b.NewBlock(proof, previousHash)
		b.NewBlock(proof, previousHash)

		// fmt.Println(fmt.Sprintf("New block mined: % +v", block.String()))
		// fmt.Println()
	}

	fmt.Println(b.String())
}
