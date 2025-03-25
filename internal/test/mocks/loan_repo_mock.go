package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

// MockLoanRepository is a mock implementation of loan.Repository
type MockLoanRepository struct {
	mock.Mock
}

// Ensure MockLoanRepository implements loan.Repository interface
var _ loan.Repository = (*MockLoanRepository)(nil)

// Get retrieves a loan by ID
func (m *MockLoanRepository) Get(ctx context.Context, id string) (*loan.Loan, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*loan.Loan), args.Error(1)
}

// Save updates an existing loan
func (m *MockLoanRepository) Save(ctx context.Context, loan *loan.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

// Create inserts a new loan
func (m *MockLoanRepository) Create(ctx context.Context, loan *loan.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

// List retrieves loans based on filter criteria
func (m *MockLoanRepository) List(ctx context.Context, filter loan.LoanFilter) ([]*loan.Loan, error) {
	args := m.Called(ctx, filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*loan.Loan), args.Error(1)
}

// Count returns the number of loans matching the filter criteria
func (m *MockLoanRepository) Count(ctx context.Context, filter loan.LoanFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// NewMockLoanRepository creates a new instance of MockLoanRepository
func NewMockLoanRepository() *MockLoanRepository {
	return &MockLoanRepository{}
}
