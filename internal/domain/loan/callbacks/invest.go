package callbacks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	loanlender "github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan_lender"
)

func (p *CallbackProvider) registerInvestCallbacks(callbacks fsm.Callbacks) {
	callbacks["before_"+loan.EventInvest] = p.beforeInvest
	callbacks["after_"+loan.EventInvest] = p.afterInvest
}

func (p *CallbackProvider) beforeInvest(ctx context.Context, e *fsm.Event) {
	loanObj := e.Args[0].(*loan.Loan)
	lender := e.Args[1].(*lender.Lender)
	amount := e.Args[2].(float64)

	if amount <= 0 {
		e.Cancel(errors.New("investment amount must be positive"))
		return
	}

	// Validate lender exists
	if _, err := p.LenderRepository.Get(ctx, lender.ID); err != nil {
		e.Cancel(errors.New("lender not found"))
		return
	}

	// Calculate current total investment
	var currentInvestment float64
	investments, err := p.LoanLenderRepository.GetByLoanID(ctx, loanObj.ID)
	if err != nil {
		e.Cancel(errors.New("error fetching investments: " + err.Error()))
		return
	}

	for _, investment := range investments {
		currentInvestment += investment.Amount
	}

	// Calculate remaining amount
	remainingPrincipal := loanObj.Amount - currentInvestment

	if amount > remainingPrincipal {
		e.Cancel(errors.New("investment exceeds remaining principal amount"))
		return
	}

	// If this will fully fund the loan, proceed with state change
	validateErr := p.Validator.Validate(loanObj, loan.Status(e.Src), loan.Status(e.Dst))
	if validateErr != nil {
		e.Cancel(validateErr)
		return
	}

}

func (p *CallbackProvider) afterInvest(ctx context.Context, e *fsm.Event) {
	loanObj := e.Args[0].(*loan.Loan)
	lender := e.Args[1].(*lender.Lender)
	amount := e.Args[2].(float64)

	// Calculate current total investment
	var currentInvestment float64
	investments, err := p.LoanLenderRepository.GetByLoanID(ctx, loanObj.ID)
	if err != nil {
		e.Cancel(errors.New("error fetching investments: " + err.Error()))
		return
	}

	for _, investment := range investments {
		currentInvestment += investment.Amount
	}

	investedTime := time.Now()

	loanLender := loanlender.LoanLender{
		ID:        uuid.New().String(),
		LoanID:    loanObj.ID,
		LenderID:  lender.ID,
		Amount:    amount,
		CreatedAt: investedTime,
	}

	createErr := p.LoanLenderRepository.Create(ctx, &loanLender)
	if createErr != nil {
		e.Cancel(errors.New("error creating investment record: " + createErr.Error()))
		return
	}

	willBeFullyFunded := (currentInvestment + amount) == loanObj.Amount

	var agreementDocLink *string

	// Update status when fully funded only
	if willBeFullyFunded {
		loanObj.Status = loan.Status(e.Dst)
		loanObj.InvestmentDate = &investedTime
		loanObj.UpdatedAt = investedTime

		loanObj.StatusTransitions = append(loanObj.StatusTransitions, loan.StatusTransition{
			From:        loan.Status(e.Src),
			To:          loan.Status(e.Dst),
			Date:        investedTime,
			Description: "Loan fully invested",
			PerformedBy: lender.ID,
		})

		// Create agreement letter
		document := &document.Document{
			ID:       uuid.New().String(),
			FileName: "agreement_" + loanObj.ID + ".pdf",
		}
		agreementDocID, err := p.DocumentRepository.Create(ctx, document)
		if err != nil {
			e.Cancel(errors.New("error creating agreement document: " + err.Error()))
			return
		}

		loanObj.AgreementDocumentID = &agreementDocID
		agreementDocLink = &document.FileName

		// Update loan in DB
		err = p.LoanRepository.Save(ctx, loanObj)
		if err != nil {
			e.Cancel(errors.New("error updating loan status: " + err.Error()))
			return
		}
	}
	if result, ok := ctx.Value(loan.InvestResultKey).(*response.LoanLenderResponse); ok {
		// Copy values to the result pointer
		*result = response.LoanLenderResponse{
			RemainingAmount:   loanObj.Amount - (currentInvestment + amount),
			InvestedAmount:    currentInvestment + amount,
			AgreementDocument: agreementDocLink,
		}
	}

}
