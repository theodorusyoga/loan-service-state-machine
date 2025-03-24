package loanlender

import (
	"context"
	"time"
)

// Repository defines the data access interface for loan-lender relationships
type Repository interface {
	Get(ctx context.Context, id string) (*LoanLender, error)
	GetByLoanID(ctx context.Context, loanID string) ([]*LoanLender, error)
	GetByLenderID(ctx context.Context, lenderID string) ([]*LoanLender, error)
	Save(ctx context.Context, loanLender *LoanLender) error
	Create(ctx context.Context, loanLender *LoanLender) error
	List(ctx context.Context, filter LoanLenderFilter) ([]*LoanLender, error)
	Count(ctx context.Context, filter LoanLenderFilter) (int64, error)
}

type LoanLenderFilter struct {
	LoanID       *string
	LenderID     *string
	MinAmount    *float64
	MaxAmount    *float64
	InvestedFrom *time.Time
	InvestedTo   *time.Time
	Page         int
	PageSize     int
}

func (f *LoanLenderFilter) WithDefaults() *LoanLenderFilter {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
	return f
}
