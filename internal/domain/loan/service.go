package loan

import (
	"context"

	"github.com/google/uuid"
	borrower "github.com/theodorusyoga/loan-service-state-machine/internal/domain/borrower"
)

// Service provides loan business operations
type Service interface {
	Create(ctx context.Context, amount float64) (*Loan, error)
	GetByID(ctx context.Context, id string) (*Loan, error)
	ChangeStatus(ctx context.Context, id string, targetStatus Status, comment, performedBy string) error
	ListLoans(ctx context.Context, filter LoanFilter) ([]*Loan, error)
}

type LoanService struct {
	repository         Repository
	borrowerRepository borrower.Repository
	validator          DefaultStatusValidator
}

func NewLoanService(r Repository, b borrower.Repository) *LoanService {
	return &LoanService{
		repository:         r,
		borrowerRepository: b,
		validator:          *NewDefaultStatusValidator(),
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, borrowerID string, amount float64, rate float64, roi float64) (*Loan, error) {
	id := uuid.New().String()

	// validate borrower ID
	if _, err := s.borrowerRepository.Get(ctx, borrowerID); err != nil {
		return nil, err
	}

	loan := NewLoan(id, borrowerID, amount, rate, roi)

	if err := s.repository.Create(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}
