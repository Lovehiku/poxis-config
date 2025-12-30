package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"paxos-lab/paxos"
)

// Shared acceptors (simulating distributed nodes)
var (
	acceptors = []*paxos.Acceptor{
		{},
		{},
		{},
	}
	mu sync.Mutex
)

// Request body structure
type ProposalRequest struct {
	ProposalNumber int    `json:"ProposalNumber"`
	Value          string `json:"Value"`
}

// HTTP handler for proposing a value
func proposeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var body ProposalRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	proposer := paxos.Proposer{
		ProposalNumber: body.ProposalNumber,
		Value:          body.Value,
		TotalAcceptors: 3, // fixed Paxos cluster size
	}

	// Lock to avoid concurrent Paxos execution
	mu.Lock()
	value := proposer.Propose(body.Value, acceptors)
	mu.Unlock()

	if value != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consensus reached: %s\n", value)
	} else {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, "Consensus not reached")
	}
}

func main() {
	http.HandleFunc("/propose", proposeHandler)

	fmt.Println("Paxos web service running on port 8080...")
	http.ListenAndServe(":8080", nil)
}