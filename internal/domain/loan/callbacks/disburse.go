package callbacks

import (
	"context"
	"errors"
	"time"

	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

func (p *CallbackProvider) registerDisburseCallbacks(callbacks fsm.Callbacks) {
	callbacks["before_"+loan.EventDisburse] = p.beforeDisburse
	callbacks["after_"+loan.EventDisburse] = p.afterDisburse
}

func (p *CallbackProvider) beforeDisburse(ctx context.Context, e *fsm.Event) {
	loanObj := e.Args[0].(*loan.Loan)
	fieldOfficerId := e.Args[1].(string)
	agreementDocFileName := e.Args[2].(string)

	if agreementDocFileName == "" {
		e.Cancel(errors.New("loan agreement document is required"))
		return
	}

	if fieldOfficerId == "" {
		e.Cancel(errors.New("field officer ID is required"))
		return
	}

	// Validate transition
	err := p.Validator.Validate(loanObj, loan.Status(e.Src), loan.Status(e.Dst))
	if err != nil {
		e.Cancel(errors.New(err.Error()))
		return
	}

	_, err = p.EmployeeRepository.Get(context.Background(), fieldOfficerId)
	if err != nil {
		e.Cancel(errors.New("field officer not found"))
		return
	}

	loanlenders, err := p.LoanLenderRepository.GetByLoanID(ctx, loanObj.ID)
	if err != nil {
		e.Cancel(err)
		return
	}
	if len(loanlenders) == 0 {
		e.Cancel(errors.New("loan must have at least one investor before disbursement"))
		return
	}

	if loanObj.ROI <= 0 {
		e.Cancel(errors.New("loan interest rate must be set before disbursement"))
		return
	}

	if loanObj.Rate <= 0 {
		e.Cancel(errors.New("loan rate must be set before disbursement"))
		return
	}
}

func (p *CallbackProvider) afterDisburse(ctx context.Context, e *fsm.Event) {
	loanObj := e.Args[0].(*loan.Loan)
	now := time.Now()
	fieldOfficerId := e.Args[1].(string)
	agreementDocFileName := e.Args[2].(string)

	agreementDoc := document.NewDocument(loanObj.ID, agreementDocFileName)
	docId, err := p.DocumentRepository.Create(context.Background(), agreementDoc)
	if err != nil {
		e.Cancel(errors.New("error saving loan agreement document"))
		return
	}

	roiAmount, err := p.calculateAndSetInvestorROI(ctx, loanObj)
	if err != nil {
		e.Cancel(err)
		return
	}

	repaymentAmount := calculateBorrowerRepayment(loanObj)

	loanObj.Status = loan.Status(e.Dst)
	loanObj.DisbursementDate = &now
	loanObj.DisbursedBy = &fieldOfficerId
	loanObj.AgreementDocumentID = &docId
	loanObj.UpdatedAt = now

	loanObj.StatusTransitions = append(loanObj.StatusTransitions, loan.StatusTransition{
		From:        loan.Status(e.Src),
		To:          loan.Status(e.Dst),
		Date:        now,
		Description: "Loan disbursed",
		PerformedBy: fieldOfficerId,
	})

	loanObj.DisbursementDate = &now

	err = p.LoanRepository.Save(context.Background(), loanObj)
	if err != nil {
		e.Cancel(errors.New("error updating loan status"))
		return
	}

	if result, ok := ctx.Value(loan.InvestResultKey).(*response.DisbursementResponse); ok {
		// Copy values to the result pointer
		*result = response.DisbursementResponse{
			AgreementDocument: &agreementDocFileName,
			DisbursementDate:  now,
			DisbursedBy:       fieldOfficerId,
			BorrowerRepayment: repaymentAmount,
			InvestorROI:       roiAmount,
		}
	}
}

func (p *CallbackProvider) calculateAndSetInvestorROI(ctx context.Context, loan *loan.Loan) (float64, error) {
	totalInvestment := 0.0

	loanLenders, err := p.LoanLenderRepository.GetByLoanID(ctx, loan.ID)
	if err != nil {
		return 0, err
	}

	for _, loanLender := range loanLenders {
		totalInvestment += loanLender.Amount
	}

	roiAmount := totalInvestment * (loan.ROI / 100)

	return totalInvestment + roiAmount, nil
}

func calculateBorrowerRepayment(loan *loan.Loan) float64 {
	interestAmount := loan.Amount * (loan.Rate / 100.0)
	totalRepayment := loan.Amount + interestAmount

	return totalRepayment
}
