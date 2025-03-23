package employee

import (
	"context"
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

func (s *EmployeeService) ListEmployees(ctx context.Context, filter EmployeeFilter) ([]*Employee, error) {
	return s.repository.List(ctx, filter)
}
