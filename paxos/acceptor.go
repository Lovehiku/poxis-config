package paxos

import "sync"

// Acceptor stores Paxos state and responds to proposer messages.
type Acceptor struct {
	mu             sync.Mutex
	promisedNumber int
	acceptedNumber int
	acceptedValue  interface{}
}

// HandlePrepare processes a Prepare request.
func (a *Acceptor) HandlePrepare(p Prepare) Promise {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Accept only higher proposal numbers
	if p.ProposalNumber > a.promisedNumber {
		a.promisedNumber = p.ProposalNumber

		return Promise{
			ProposalNumber: p.ProposalNumber,
			AcceptedValue:  a.acceptedValue,
		}
	}

	// Reject by returning empty promise
	return Promise{}
}

// HandleAccept processes an Accept request.
func (a *Acceptor) HandleAccept(ac Accept) Accepted {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Accept if proposal number is valid
	if ac.ProposalNumber >= a.promisedNumber {
		a.promisedNumber = ac.ProposalNumber
		a.acceptedNumber = ac.ProposalNumber
		a.acceptedValue = ac.Value

		return Accepted{
			ProposalNumber: ac.ProposalNumber,
			Value:          ac.Value,
		}
	}

	// Reject by returning empty Accepted
	return Accepted{}
}
