package borrower

import (
	"context"
)

// Service provides borrower business operations
type Service interface {
	Create(ctx context.Context, fullName, email, phoneNumber, idNumber string) (*Borrower, error)
	ListBorrowers(ctx context.Context, filter BorrowerFilter) ([]*Borrower, error)
}

type BorrowerService struct {
	repository Repository
}

func NewBorrowerService(r Repository) *BorrowerService {
	return &BorrowerService{
		repository: r,
	}
}

func (s *BorrowerService) CreateBorrower(ctx context.Context, fullName, email, phoneNumber, idNumber string) (*Borrower, error) {
	borrower := NewBorrower(fullName, email, phoneNumber, idNumber)

	if err := s.repository.Create(ctx, borrower); err != nil {
		return nil, err
	}

	return borrower, nil
}

func (s *BorrowerService) GetByID(ctx context.Context, id string) (*Borrower, error) {
	return s.repository.Get(ctx, id)
}
