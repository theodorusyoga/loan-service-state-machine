package callbacks

import (
	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

// CallbackProvider provides callback functions for the loan state machine
type CallbackProvider struct {
	LoanRepository     loan.Repository
	EmployeeRepository employee.Repository
	DocumentRepository document.Repository
	Validator          loan.DefaultStatusValidator
}

func New(
	loanRepo loan.Repository,
	empRepo employee.Repository,
	docRepo document.Repository,
) *CallbackProvider {
	return &CallbackProvider{
		LoanRepository:     loanRepo,
		EmployeeRepository: empRepo,
		DocumentRepository: docRepo,
		Validator:          *loan.NewDefaultStatusValidator(),
	}
}

func (p *CallbackProvider) GetCallbacks() fsm.Callbacks {
	callbacks := fsm.Callbacks{}

	// Add approve callbacks
	p.registerApproveCallbacks(callbacks)

	// TODO: Add other callbacks

	return callbacks
}
