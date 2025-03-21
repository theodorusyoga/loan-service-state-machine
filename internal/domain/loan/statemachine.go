package loan

import (
	"context"
	"errors"
	"time"

	"github.com/looplab/fsm"
)

const (
	EventApprove  = "approve"
	EventInvest   = "invest"
	EventDisburse = "disburse"
	EventReject   = "reject"
)

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
		fsm.Callbacks{
			"before_" + EventApprove: func(_ context.Context, e *fsm.Event) {
				// TODO: Check document completeness
			},
			"after_" + EventApprove: func(_ context.Context, e *fsm.Event) {
				loan := e.Args[0].(*Loan)
				now := time.Now()
				approvedBy := e.Args[1].(string)

				loan.Status = Status(e.Dst)
				loan.ApprovalDate = &now
				loan.ApprovedBy = approvedBy
				loan.UpdatedAt = now

				loan.StatusTransitions = append(loan.StatusTransitions, StatusTransition{
					From:        Status(e.Src),
					To:          Status(e.Dst),
					Date:        now,
					Description: "Loan approved",
					PerformedBy: approvedBy,
				})
			},
		},
	)
}

func (s *LoanService) ApproveLoan(loan *Loan, approvedBy string) error {
	loanFSM := s.createFSM(loan)
	err := loanFSM.Event(context.Background(), EventApprove, loan, approvedBy)
	if err != nil {
		if errors.Is(err, fsm.NoTransitionError{}) {
			return errors.New("cannot approve loan in current state")
		}
		return err
	}
	return nil
}
