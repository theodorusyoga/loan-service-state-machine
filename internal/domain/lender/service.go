package lender

import (
	"context"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain"
)

// Service provides lender business operations
type Service interface {
	Create(ctx context.Context, fullName, email, phoneNumber, idNumber string) (*Lender, error)
	ListLenders(ctx context.Context, filter LenderFilter) ([]*Lender, error)
}

type LenderService struct {
	repository Repository
}

func NewLenderService(r Repository) *LenderService {
	return &LenderService{
		repository: r,
	}
}

func (s *LenderService) CreateLender(ctx context.Context, fullName, email, phoneNumber, idNumber string) (*Lender, error) {
	lender := NewLender(fullName, email, phoneNumber, idNumber)

	if err := s.repository.Create(ctx, lender); err != nil {
		return nil, err
	}

	return lender, nil
}

func (s *LenderService) GetByID(ctx context.Context, id string) (*Lender, error) {
	return s.repository.Get(ctx, id)
}

func (s *LenderService) ListLenders(ctx context.Context, filter LenderFilter) (*domain.PaginatedResponse, error) {
	filter.WithDefaults()
	lenders, err := s.repository.List(ctx, filter)
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
		Data: lenders,
		Pagination: domain.PaginationInfo{
			CurrentPage: filter.Page,
			PageSize:    filter.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}
