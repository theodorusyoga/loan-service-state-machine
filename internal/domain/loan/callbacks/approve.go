package callbacks

import (
	"context"
	"errors"
	"time"

	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

func (p *CallbackProvider) registerApproveCallbacks(callbacks fsm.Callbacks) {
	callbacks["before_"+loan.EventApprove] = p.beforeApprove
	callbacks["after_"+loan.EventApprove] = p.afterApprove
}

func (p *CallbackProvider) beforeApprove(ctx context.Context, e *fsm.Event) {
	// TODO: Check document completeness
	loanObj := e.Args[0].(*loan.Loan)
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
	err := p.Validator.Validate(loanObj, loan.Status(e.Src), loan.Status(e.Dst))
	if err != nil {
		e.Cancel(errors.New(err.Error()))
		return
	}

	// check employee exists in DB
	_, err = p.EmployeeRepository.Get(context.Background(), approvedBy)
	if err != nil {
		e.Cancel(errors.New("employee not found"))
		return
	}
}

func (p *CallbackProvider) afterApprove(ctx context.Context, e *fsm.Event) {
	loanObj := e.Args[0].(*loan.Loan)
	now := time.Now()
	approvedBy := e.Args[1].(string)

	// validate document ID exists
	fileName := e.Args[2].(string)

	// insert document
	doc := document.NewDocument(loanObj.ID, fileName)
	// create document
	docId, docErr := p.DocumentRepository.Create(context.Background(), doc)
	if docErr != nil {
		e.Cancel(errors.New("error creating document"))
		return
	}

	loanObj.Status = loan.Status(e.Dst)
	loanObj.ApprovalDate = &now
	loanObj.ApprovedBy = &approvedBy
	loanObj.SurveyDocumentID = &docId
	loanObj.UpdatedAt = now

	loanObj.StatusTransitions = append(loanObj.StatusTransitions, loan.StatusTransition{
		From:        loan.Status(e.Src),
		To:          loan.Status(e.Dst),
		Date:        now,
		Description: "Loan approved",
		PerformedBy: approvedBy,
	})

	// update to DB
	// TODO: Should be in transaction
	err := p.LoanRepository.Save(context.Background(), loanObj)
	if err != nil {
		e.Cancel(errors.New("error updating loan status"))
		return
	}

}
