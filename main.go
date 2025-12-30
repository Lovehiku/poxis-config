package main

import (
	"fmt"
	"paxos-lab/paxos"
)

func main() {
	// Create 3 acceptors (2k+1)
	acceptors := []*paxos.Acceptor{
		{},
		{},
		
		
		
	}

	// Create proposer
	proposer := paxos.Proposer{
		ProposalNumber: 1,
		Value:          "Distributed Systems",
		TotalAcceptors: 5,
	}

	// Run Paxos
	value := proposer.Propose("Distributed Systems", acceptors)

	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
	}
}
