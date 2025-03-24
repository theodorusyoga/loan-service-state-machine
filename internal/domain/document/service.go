package document

import (
	"context"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain"
)

type DocumentService struct {
	repository Repository
}

func NewDocumentService(r Repository) *DocumentService {
	return &DocumentService{
		repository: r,
	}
}

func (s *DocumentService) CreateDocument(ctx context.Context, loanID, fileName string) (*Document, error) {
	document := NewDocument(loanID, fileName)

	if _, err := s.repository.Create(ctx, document); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *DocumentService) GetByID(ctx context.Context, id string) (*Document, error) {
	return s.repository.Get(ctx, id)
}

func (s *DocumentService) ListDocuments(ctx context.Context, filter DocumentFilter) (*domain.PaginatedResponse, error) {
	filter.WithDefaults()
	documents, err := s.repository.List(ctx, filter)
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
		Data: documents,
		Pagination: domain.PaginationInfo{
			CurrentPage: filter.Page,
			PageSize:    filter.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}
