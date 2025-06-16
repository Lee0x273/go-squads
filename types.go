package squads

import "squads/generated/squads_multisig_program"

type Permission uint8

const (
	Initiate Permission = 1 << 0
	Vote     Permission = 1 << 1
	Execute  Permission = 1 << 2
)

func (p Permission) Has(permission Permission) bool {
	return p&permission != 0
}

type VoteOP uint8

const (
	VoteOPApprove VoteOP = iota
	VoteOPReject
	VoteOPCancel
)

type ProposalStatus uint8

const (
	ProposalStatusDraft ProposalStatus = iota
	ProposalStatusActive
	ProposalStatusRejected
	ProposalStatusApproved
	ProposalStatusExecuting
	ProposalStatusExecuted
	ProposalStatusCancelled
)

func GetProposalStatus(status squads_multisig_program.ProposalStatus) ProposalStatus {
	switch status.(type) {
	case *squads_multisig_program.ProposalStatusDraft:
		return ProposalStatusDraft
	case *squads_multisig_program.ProposalStatusActive:
		return ProposalStatusActive
	case *squads_multisig_program.ProposalStatusRejected:
		return ProposalStatusRejected
	case *squads_multisig_program.ProposalStatusApproved:
		return ProposalStatusApproved
	case *squads_multisig_program.ProposalStatusExecuting:
		return ProposalStatusExecuting
	case *squads_multisig_program.ProposalStatusExecuted:
		return ProposalStatusExecuted
	case *squads_multisig_program.ProposalStatusCancelled:
		return ProposalStatusCancelled
	}
	return ProposalStatusDraft
}
