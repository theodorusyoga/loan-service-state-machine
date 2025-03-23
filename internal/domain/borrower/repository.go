package borrower

import (
	"context"
)

// Repository defines the data access interface for borrowers
type Repository interface {
	Get(ctx context.Context, id string) (*Borrower, error)
	Save(ctx context.Context, borrower *Borrower) error
	Create(ctx context.Context, borrower *Borrower) error
	// TODO: Implement the following methods
	List(ctx context.Context, filter BorrowerFilter) ([]*Borrower, error)
	// Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter BorrowerFilter) (int64, error)
}

type BorrowerFilter struct {
	FullName    *string
	Email       *string
	PhoneNumber *string
	IDNumber    *string
	Page        int
	PageSize    int
}

func (f *BorrowerFilter) WithDefaults() *BorrowerFilter {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
	return f
}
