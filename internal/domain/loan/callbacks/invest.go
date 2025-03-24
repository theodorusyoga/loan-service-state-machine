package callbacks

import (
	"context"
	"errors"
	"time"

	"github.com/looplab/fsm"
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
	lender := e.Args[1].(lender.Lender)
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
	lender := e.Args[1].(lender.Lender)
	amount := e.Args[2].(float64)
	performedBy := e.Args[3].(string)

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
			PerformedBy: performedBy,
		})

		// Update loan in DB
		err = p.LoanRepository.Save(ctx, loanObj)
		if err != nil {
			e.Cancel(errors.New("error updating loan status: " + err.Error()))
			return
		}
	}
}
