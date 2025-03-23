package employee

import (
	"context"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain"
)

// Service provides employee business operations
type Service interface {
	CreateEmployee(ctx context.Context, fullName, email, phoneNumber, idNumber string) (*Employee, error)
	GetByID(ctx context.Context, id string) (*Employee, error)
	ListEmployees(ctx context.Context, filter EmployeeFilter) ([]*Employee, error)
}

type EmployeeService struct {
	repository Repository
}

func NewEmployeeService(r Repository) *EmployeeService {
	return &EmployeeService{
		repository: r,
	}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, fullName, email, phoneNumber, idNumber string) (*Employee, error) {
	employee := NewEmployee(fullName, email, phoneNumber, idNumber)

	if err := s.repository.Create(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *EmployeeService) GetByID(ctx context.Context, id string) (*Employee, error) {
	return s.repository.Get(ctx, id)
}

func (s *EmployeeService) ListEmployees(ctx context.Context, filter EmployeeFilter) (*domain.PaginatedResponse, error) {
	filter.WithDefaults()
	employees, err := s.repository.List(ctx, filter)
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
		Data: employees,
		Pagination: domain.PaginationInfo{
			CurrentPage: filter.Page,
			PageSize:    filter.PageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}
