package employee

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents a domain entity for an employee
type Employee struct {
	ID          string
	FullName    string
	Email       string
	PhoneNumber string
	IDNumber    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewEmployee(fullName, email, phoneNumber, idNumber string) *Employee {
	now := time.Now()
	return &Employee{
		ID:          uuid.New().String(),
		FullName:    fullName,
		Email:       email,
		PhoneNumber: phoneNumber,
		IDNumber:    idNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
