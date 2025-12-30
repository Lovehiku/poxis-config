package paxos

// Prepare is sent by the proposer to ask acceptors
// if they will promise to a proposal number.
type Prepare struct {
	ProposalNumber int
}

// Promise is sent by acceptors in response to Prepare.
// It may include a previously accepted value.
type Promise struct {
	ProposalNumber int
	AcceptedValue  interface{}
}

// Accept is sent by the proposer asking acceptors
// to accept a specific value.
type Accept struct {
	ProposalNumber int
	Value          interface{}
}

// Accepted is sent by acceptors confirming acceptance.
type Accepted struct {
	ProposalNumber int
	Value          interface{}
}
