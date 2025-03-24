package loan

import (
	"context"
	"errors"
	"time"

	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
)

const (
	EventApprove  = "approve"
	EventInvest   = "invest"
	EventDisburse = "disburse"
	EventReject   = "reject"
)

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
		fsm.Callbacks{
			"before_" + EventApprove: func(_ context.Context, e *fsm.Event) {
				// TODO: Check document completeness
				loan := e.Args[0].(*Loan)
				approvedBy := e.Args[1].(string)
				fileName := e.Args[2].(string)

				if fileName == "" {
					e.Cancel(errors.New("document is required"))
					return
				}

				if approvedBy == "" {
					e.Cancel(errors.New("approved by is required"))
					return
				}

				// validate transition
				err := s.validator.Validate(loan, Status(e.Src), Status(e.Dst))
				if err != nil {
					e.Cancel(errors.New(err.Error()))
					return
				}

				// check employee exists in DB
				_, err = s.employeeRepository.Get(context.Background(), approvedBy)
				if err != nil {
					e.Cancel(errors.New("employee not found"))
					return
				}

			},
			"after_" + EventApprove: func(_ context.Context, e *fsm.Event) {
				loan := e.Args[0].(*Loan)
				now := time.Now()
				approvedBy := e.Args[1].(string)

				// validate document ID exists
				fileName := e.Args[2].(string)

				// insert document
				doc := document.NewDocument(loan.ID, fileName)
				// create document
				docId, docErr := s.documentRepository.Create(context.Background(), doc)
				if docErr != nil {
					e.Cancel(errors.New("error creating document"))
					return
				}

				loan.Status = Status(e.Dst)
				loan.ApprovalDate = &now
				loan.ApprovedBy = &approvedBy
				loan.SurveyDocumentID = &docId
				loan.UpdatedAt = now

				loan.StatusTransitions = append(loan.StatusTransitions, StatusTransition{
					From:        Status(e.Src),
					To:          Status(e.Dst),
					Date:        now,
					Description: "Loan approved",
					PerformedBy: approvedBy,
				})

				// update to DB
				// TODO: Should be in transaction
				err := s.repository.Save(context.Background(), loan)
				if err != nil {
					e.Cancel(errors.New("error updating loan status"))
					return
				}

			},
		},
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
