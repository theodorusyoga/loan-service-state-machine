package loan

import (
	"context"
)

// Repository defines the data access interface for loans
type Repository interface {
	Get(ctx context.Context, id string) (*Loan, error)
	Save(ctx context.Context, loan *Loan) error
	Create(ctx context.Context, loan *Loan) error
	// TODO: Implement the following methods
	List(ctx context.Context, filter LoanFilter) ([]*Loan, error)
	// Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter LoanFilter) (int64, error)
}

type LoanFilter struct {
	Status    *Status
	MinAmount *float64
	MaxAmount *float64
	Page      int
	PageSize  int
}

func (f *LoanFilter) WithDefaults() *LoanFilter {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
	return f
}
