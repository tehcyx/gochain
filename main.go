package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tehcyx/gochain/blockchain"
)

var b = blockchain.NewBlockchain()
var nodeIdentifier uuid.UUID

type chainResponse struct {
	Chain  []blockchain.Block `json:"chain"`
	Length int                `json:"length"`
}

type nodeRequest struct {
	Nodes []node `json:"nodes"`
}
type node struct {
	Address string `json:"address"`
	Comment string `json:"comment"`
}

func main() {
	if len(b.Chain) != 1 {
		log.Fatal("Error occurred initializing the chain")
	}

	var err error
	nodeIdentifier, err = uuid.NewRandom()
	if err != nil {
		log.Fatal("Could not create node identifier")
	}
	fmt.Println(fmt.Sprintf("Node Identifier is %s", nodeIdentifier.String()))

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	//r.HandleFunc("/transactions", transactionsHandler).Methods("GET")
	r.HandleFunc("/transactions/new", transactionsNewHandler).Methods("POST")
	r.HandleFunc("/mine", mineHandler).Methods("GET")
	r.HandleFunc("/chain", chainHandler).Methods("GET")
	r.HandleFunc("/nodes/register", nodeRegisterHandler).Methods("POST")
	r.HandleFunc("/nodes/resolve", nodeResolveHandler).Methods("GET")
	http.Handle("/", r)

	port := ":8080"

	srv := &http.Server{
		Handler: r,
		Addr:    port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println(fmt.Sprintf("Listening on %s", port))
	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok index"))
}

func mineHandler(w http.ResponseWriter, r *http.Request) {
	lastBlock := b.LastBlock()
	lastProof := lastBlock.Proof
	proof := b.ProofOfWork(lastProof)

	fmt.Println(fmt.Sprintf("Transaction for mining reward of Block %d", b.NewTransaction("0", nodeIdentifier.String(), 1.0)))

	previousHash := b.Hash(lastBlock)
	block := b.NewBlock(proof, previousHash)

	var response = map[string]interface{}{
		"message":     fmt.Sprintf("Success: New block mined %d", block.Index),
		"block":       block,
		"status_code": http.StatusCreated,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok transactions"))
}

func transactionsNewHandler(w http.ResponseWriter, r *http.Request) {
	var t blockchain.Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		fmt.Println(fmt.Sprintf("Error decoding request payload: %s", err))
		var response = map[string]interface{}{
			"message":     "Error: Invalid request payload",
			"status_code": http.StatusBadRequest,
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	// Check that the required fields are in the POST'ed data
	if !(t.Sender != "" && t.Recipient != "" && t.Amount > 0.0) {
		var response = map[string]interface{}{
			"message":     "Error: Invalid request payload",
			"status_code": http.StatusBadRequest,
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}
	blockNumber := b.NewTransaction(t.Sender, t.Recipient, t.Amount)
	var response = map[string]interface{}{
		"message":     fmt.Sprintf("Success: Transaction will be added to block %d", blockNumber),
		"status_code": http.StatusCreated,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func chainHandler(w http.ResponseWriter, r *http.Request) {
	var chainJSONResponse chainResponse
	chainJSONResponse.Chain = b.Chain
	chainJSONResponse.Length = len(b.Chain)

	respondWithJSON(w, http.StatusOK, chainJSONResponse)
}

func nodeRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var nodes nodeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&nodes); err != nil {
		fmt.Println(fmt.Sprintf("Error decoding request payload: %s", err))
		var response = map[string]interface{}{
			"message":     "Error: Invalid request payload",
			"status_code": http.StatusBadRequest,
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}
	if nodes.Nodes == nil || len(nodes.Nodes) == 0 {
		var response = map[string]interface{}{
			"message":     "Error: Please provide a valid list of nodes",
			"status_code": http.StatusBadRequest,
		}
		respondWithJSON(w, http.StatusBadRequest, response)
		return
	}
	var success []bool
	for nodeInfo := range nodes.Nodes {
		isAdded := b.RegisterNode(nodes.Nodes[nodeInfo].Address, nodes.Nodes[nodeInfo].Comment)
		if isAdded {
			success = append(success, true)
		}
	}
	var response = map[string]interface{}{
		"message":     fmt.Sprintf("Success: %d node(s) successfully added", len(success)),
		"nodes":       b.Nodes,
		"status_code": http.StatusCreated,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func nodeResolveHandler(w http.ResponseWriter, r *http.Request) {
	replaced := resolveConflicts(b)

	if replaced {
		var response = map[string]interface{}{
			"message":     "Our chain was replaced",
			"chain":       b.Chain,
			"status_code": http.StatusOK,
		}
		respondWithJSON(w, http.StatusOK, response)
		return
	}
	var response = map[string]interface{}{
		"message":     "Our chain is authoritative",
		"chain":       b.Chain,
		"status_code": http.StatusOK,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func validChain(self blockchain.Blockchain, other []blockchain.Block) bool {
	lastBlock := other[0]
	currentIndex := 1

	for currentIndex < len(other) {
		block := other[currentIndex]
		if block.PreviousHash != self.Hash(lastBlock) {
			return false
		}
		if !self.ValidProof(lastBlock.Proof, block.Proof) {
			return false
		}

		lastBlock = block
		currentIndex++
	}

	return true
}

func resolveConflicts(self *blockchain.Blockchain) bool {
	neighbours := self.Nodes
	var newChain []blockchain.Block

	maxLength := len(self.Chain)

	for node := range neighbours {
		fmt.Println("http://" + node + "/chain")
		resp, err := http.Get("http://" + node + "/chain")
		if err != nil {
			log.Fatal(fmt.Sprintf("Error calling neighbours chain: %s", err))
		}
		if resp.StatusCode == 200 {
			var cr chainResponse
			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(&cr); err != nil {
				fmt.Println(fmt.Sprintf("Error decoding request payload: %s", err))
				return false
			}
			if cr.Length > maxLength && validChain(*self, cr.Chain) {
				maxLength = cr.Length
				newChain = cr.Chain
			}
		}
	}
	if newChain != nil {
		self.Chain = newChain
		return true
	}
	return false
}
