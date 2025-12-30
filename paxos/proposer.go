package paxos

// Proposer initiates Paxos consensus.
type Proposer struct {
	ProposalNumber int
	Value          interface{}
	TotalAcceptors int
}

// Propose runs the Paxos algorithm.
func (p *Proposer) Propose(value interface{}, acceptors []*Acceptor) interface{} {
	promises := 0

	// Phase 1: Prepare
	for _, acceptor := range acceptors {
		promise := acceptor.HandlePrepare(
			Prepare{ProposalNumber: p.ProposalNumber},
		)

		if promise.ProposalNumber == p.ProposalNumber {
			promises++
		}
	}

	// Check for majority promises
	if promises > p.TotalAcceptors/2 {
		accepted := 0

		// Phase 2: Accept
		for _, acceptor := range acceptors {
			ack := acceptor.HandleAccept(
				Accept{
					ProposalNumber: p.ProposalNumber,
					Value:          value,
				},
			)

			if ack.ProposalNumber == p.ProposalNumber {
				accepted++
			}
		}

		// Check for majority acceptance
		if accepted > len(acceptors)/2 {
			return value
		}
	}

	// Consensus failed
	return nil
}
