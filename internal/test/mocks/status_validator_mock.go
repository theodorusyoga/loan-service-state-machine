package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

// MockStatusValidator is a mock implementation of loan.StatusValidator
type MockStatusValidator struct {
	mock.Mock
}

// Ensure MockStatusValidator implements loan.StatusValidator interface
var _ loan.StatusValidator = (*MockStatusValidator)(nil)

func (m *MockStatusValidator) Validate(loan *loan.Loan, from, to loan.Status) error {
	args := m.Called(loan, from, to)
	return args.Error(0)
}

func NewMockStatusValidator() *MockStatusValidator {
	return &MockStatusValidator{}
}
