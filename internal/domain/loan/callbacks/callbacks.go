package callbacks

import (
	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	loanlender "github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan_lender"
)

// CallbackProvider provides callback functions for the loan state machine
type CallbackProvider struct {
	LenderRepository     loanlender.Repository
	LoanRepository       loan.Repository
	LoanLenderRepository loanlender.Repository
	EmployeeRepository   employee.Repository
	DocumentRepository   document.Repository
	Validator            loan.DefaultStatusValidator
}

func New(
	lenderRepo loanlender.Repository,
	loanRepo loan.Repository,
	loanLenderRepo loanlender.Repository,
	empRepo employee.Repository,
	docRepo document.Repository,
) *CallbackProvider {
	return &CallbackProvider{
		LenderRepository:     lenderRepo,
		LoanRepository:       loanRepo,
		LoanLenderRepository: loanLenderRepo,
		EmployeeRepository:   empRepo,
		DocumentRepository:   docRepo,
		Validator:            *loan.NewDefaultStatusValidator(),
	}
}

func (p *CallbackProvider) GetCallbacks() fsm.Callbacks {
	callbacks := fsm.Callbacks{}

	// Add approve callbacks
	p.registerApproveCallbacks(callbacks)

	// Add invest callbacks
	p.registerInvestCallbacks(callbacks)

	// TODO: Add other callbacks

	return callbacks
}
