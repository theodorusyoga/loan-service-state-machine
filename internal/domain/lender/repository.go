package lender

import (
	"context"
)

// Repository defines the data access interface for borrowers
type Repository interface {
	Get(ctx context.Context, id string) (*Lender, error)
	Save(ctx context.Context, borrower *Lender) error
	Create(ctx context.Context, borrower *Lender) error
	// TODO: Implement the following methods
	List(ctx context.Context, filter LenderFilter) ([]*Lender, error)
	// Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter LenderFilter) (int64, error)
}

type LenderFilter struct {
	FullName    *string
	Email       *string
	PhoneNumber *string
	IDNumber    *string
	Page        int
	PageSize    int
}

func (f *LenderFilter) WithDefaults() *LenderFilter {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
	return f
}
