package callbacks

import (
	"context"

	"github.com/looplab/fsm"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	loanlender "github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan_lender"
)

type LoanCallbackProvider interface {
	// GetCallbacks returns all registered callbacks for the FSM
	GetCallbacks() fsm.Callbacks

	// Specific callback methods for direct testing
	BeforeApproval(ctx context.Context, e *fsm.Event)
	AfterApproval(ctx context.Context, e *fsm.Event)

	BeforeInvest(ctx context.Context, e *fsm.Event)
	AfterInvest(ctx context.Context, e *fsm.Event)

	BeforeDisburse(ctx context.Context, e *fsm.Event)
	AfterDisburse(ctx context.Context, e *fsm.Event)
}

// CallbackProvider provides callback functions for the loan state machine
type CallbackProvider struct {
	LenderRepository     lender.Repository
	LoanRepository       loan.Repository
	LoanLenderRepository loanlender.Repository
	EmployeeRepository   employee.Repository
	DocumentRepository   document.Repository
	Validator            loan.DefaultStatusValidator
}

var _ LoanCallbackProvider = (*CallbackProvider)(nil)

func New(
	lenderRepo lender.Repository,
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

	// Add disburse callbacks
	p.registerDisburseCallbacks(callbacks)

	return callbacks
}
