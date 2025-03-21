package loan

import (
	"context"

	"github.com/google/uuid"
)

// Service provides loan business operations
type Service interface {
	Create(ctx context.Context, amount float64) (*Loan, error)
	GetByID(ctx context.Context, id string) (*Loan, error)
	ChangeStatus(ctx context.Context, id string, targetStatus Status, comment, performedBy string) error
	ListLoans(ctx context.Context, filter LoanFilter) ([]*Loan, error)
}

type LoanService struct {
	repository Repository
	validator  DefaultStatusValidator
}

func NewLoanService(r Repository) *LoanService {
	return &LoanService{
		repository: r,
		validator:  *NewDefaultStatusValidator(),
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, amount float64) (*Loan, error) {
	id := uuid.New().String()
	loan := NewLoan(id, amount)

	if err := s.repository.Save(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}
