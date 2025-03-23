package borrower

import (
	"context"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain"
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

func (s *BorrowerService) ListBorrowers(ctx context.Context, filter BorrowerFilter) (*domain.PaginatedResponse, error) {
	filter.WithDefaults()
	borrowers, err := s.repository.List(ctx, filter)
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
		Data: borrowers,
		Pagination: domain.PaginationInfo{
			CurrentPage: filter.Page,
			PageSize:    filter.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}
