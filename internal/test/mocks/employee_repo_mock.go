package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
)

// MockEmployeeRepository is a mock implementation of employee.Repository
type MockEmployeeRepository struct {
	mock.Mock
}

// Ensure MockEmployeeRepository implements employee.Repository interface
var _ employee.Repository = (*MockEmployeeRepository)(nil)

// Get retrieves an employee by ID
func (m *MockEmployeeRepository) Get(ctx context.Context, id string) (*employee.Employee, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	// return mock Employee
	return employee.NewEmployee("", "", "", ""), args.Error(1)
}

// Save updates an existing employee
func (m *MockEmployeeRepository) Save(ctx context.Context, emp *employee.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

// Create inserts a new employee
func (m *MockEmployeeRepository) Create(ctx context.Context, emp *employee.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

// List retrieves employees based on filter criteria
func (m *MockEmployeeRepository) List(ctx context.Context, filter employee.EmployeeFilter) ([]*employee.Employee, error) {
	args := m.Called(ctx, filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*employee.Employee), args.Error(1)
}

// Delete removes an employee by ID
func (m *MockEmployeeRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// NewMockEmployeeRepository creates a new instance of MockEmployeeRepository
func NewMockEmployeeRepository() *MockEmployeeRepository {
	return &MockEmployeeRepository{}
}

// Count implements employee.Repository.
func (m *MockEmployeeRepository) Count(ctx context.Context, filter employee.EmployeeFilter) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
