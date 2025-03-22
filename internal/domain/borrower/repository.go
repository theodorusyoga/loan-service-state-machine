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
	// List(ctx context.Context, filter BorrowerFilter) ([]*Borrower, error)
	// Delete(ctx context.Context, id string) error
}

type BorrowerFilter struct {
	FullName    *string
	Email       *string
	PhoneNumber *string
	IDNumber    *string
	Page        int
	PageSize    int
}
