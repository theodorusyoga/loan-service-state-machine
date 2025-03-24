package loan

import (
	"context"
	"errors"

	"github.com/looplab/fsm"
)

const (
	EventApprove  = "approve"
	EventInvest   = "invest"
	EventDisburse = "disburse"
	EventReject   = "reject"
)

type CallbackRegistrar interface {
	GetCallbacks() fsm.Callbacks
}

func IsValidStatus(status string) bool {
	validStatuses := []string{
		string(EventApprove),
		string(EventDisburse),
		string(EventInvest),
		string(EventReject),
	}

	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}

	return false
}

// Create finite state machine for loan status
func (s *LoanService) createFSM(loan *Loan) *fsm.FSM {
	return fsm.NewFSM(
		string(loan.Status),
		fsm.Events{
			{Name: EventApprove, Src: []string{string(StatusProposed)}, Dst: string(StatusApproved)},
			{Name: EventInvest, Src: []string{string(StatusApproved)}, Dst: string(StatusInvested)},
			{Name: EventDisburse, Src: []string{string(StatusInvested)}, Dst: string(StatusDisbursed)},
			{Name: EventReject, Src: []string{string(StatusProposed)}},
		},
		s.callbackRegistrar.GetCallbacks(),
	)
}

func (s *LoanService) ApproveLoan(loan *Loan, approvedBy string, fileName string) error {
	loanFSM := s.createFSM(loan)
	err := loanFSM.Event(context.Background(), EventApprove, loan, approvedBy, fileName)
	if err != nil {
		if errors.Is(err, fsm.NoTransitionError{}) {
			return errors.New("cannot approve loan in current state")
		}
		return err
	}
	return nil
}
