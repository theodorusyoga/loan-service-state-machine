package employee

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents a domain entity for an employee
type Employee struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IDNumber    string    `json:"id_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
