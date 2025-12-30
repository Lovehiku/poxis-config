package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "sync"
    "time"
)

// ---------------- Paxos Simulation ----------------

// Acceptor represents a node that can accept or reject proposals
type Acceptor struct {
    mu       sync.Mutex
    promised int
    accepted interface{}
}

func (a *Acceptor) Accept(proposalNumber int, value interface{}) bool {
    a.mu.Lock()
    defer a.mu.Unlock()

    // Reject if proposal number is lower than promised
    if proposalNumber < a.promised {
        return false
    }

    // Randomly reject to simulate network unreliability
    if rand.Float32() < 0.3 {
        return false
    }

    a.promised = proposalNumber
    a.accepted = value
    return true
}

// Proposer tries to get consensus among acceptors
type Proposer struct{}

func (p *Proposer) Propose(proposalNumber int, value interface{}, acceptors []*Acceptor) interface{} {
    acceptedCount := 0
    for _, a := range acceptors {
        if a.Accept(proposalNumber, value) {
            acceptedCount++
        }
    }
    // Consensus requires majority
    if acceptedCount > len(acceptors)/2 {
        return value
    }
    return nil
}

// ---------------- Web Service ----------------

// Shared acceptors
var (
    acceptors = []*Acceptor{
        &Acceptor{},
        &Acceptor{},
        &Acceptor{},
    }
    proposer   = &Proposer{}
    maxRetries = 3
)

// Request body
type ProposalRequest struct {
    ProposalNumber int    `json:"ProposalNumber"`
    Value          string `json:"Value"`
}

// HTTP handler
func proposeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }

    var body ProposalRequest
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    resultChan := make(chan interface{})

    go func() {
        var result interface{}
        for attempt := 1; attempt <= maxRetries; attempt++ {
            log.Printf("Paxos attempt %d\n", attempt)
            result = proposer.Propose(body.ProposalNumber, body.Value, acceptors)
            if result != nil {
                break
            }
            time.Sleep(200 * time.Millisecond)
        }
        resultChan <- result
    }()

    select {
    case result := <-resultChan:
        if result != nil {
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "Consensus reached: %s\n", result)
        } else {
            w.WriteHeader(http.StatusConflict)
            fmt.Fprintln(w, "Consensus not reached")
        }
    case <-ctx.Done():
        log.Println("Paxos timeout occurred")
        http.Error(w, "Paxos timeout", http.StatusRequestTimeout)
    }
}

func main() {
    rand.Seed(time.Now().UnixNano()) // randomness for rejection
    http.HandleFunc("/propose", proposeHandler)
    fmt.Println("Paxos web service running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
