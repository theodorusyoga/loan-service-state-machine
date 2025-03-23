package loan

import (
	"context"

	"github.com/google/uuid"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain"
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

func (s *LoanService) GetByID(ctx context.Context, id string) (*Loan, error) {
	return s.repository.Get(ctx, id)
}

func (s *LoanService) ListLoans(ctx context.Context, filter LoanFilter) (*domain.PaginatedResponse, error) {
	filter.WithDefaults()
	loans, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get the total count
	totalItems, err := s.repository.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 0
	if filter.PageSize > 0 {
		totalPages = int((totalItems + int64(filter.PageSize) - 1) / int64(filter.PageSize))
	}

	return &domain.PaginatedResponse{
		Data: loans,
		Pagination: domain.PaginationInfo{
			CurrentPage: filter.Page,
			PageSize:    filter.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}
