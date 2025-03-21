package loan

import (
	"context"
)

// Repository defines the data access interface for loans
type Repository interface {
	Get(ctx context.Context, id string) (*Loan, error)
	Save(ctx context.Context, loan *Loan) error
	List(ctx context.Context, filter LoanFilter) ([]*Loan, error)
	Delete(ctx context.Context, id string) error
}

type LoanFilter struct {
	Status    *Status
	MinAmount *float64
	MaxAmount *float64
	Page      int
	PageSize  int
}
