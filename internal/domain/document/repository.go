package document

import (
	"context"
)

// Repository defines the data access interface for documents
type Repository interface {
	Get(ctx context.Context, id string) (*Document, error)
	Save(ctx context.Context, document *Document) error
	Create(ctx context.Context, document *Document) error
	// TODO: Implement the following methods
	List(ctx context.Context, filter DocumentFilter) ([]*Document, error)
	// Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter DocumentFilter) (int64, error)
}

type DocumentFilter struct {
	LoanID   *string
	FileName *string
	Page     int
	PageSize int
}

func (f *DocumentFilter) WithDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
}
