package loan

import (
	"context"
)

// Service provides loan business operations
type Service interface {
	Create(ctx context.Context, amount float64) (*Loan, error)
	GetByID(ctx context.Context, id string) (*Loan, error)
	ChangeStatus(ctx context.Context, id string, targetStatus Status, comment, performedBy string) error
	ListLoans(ctx context.Context, filter LoanFilter) ([]*Loan, error)
}

type LoanService struct {
	Repository Repository
}
