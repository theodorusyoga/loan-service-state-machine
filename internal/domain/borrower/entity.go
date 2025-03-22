package borrower

import (
	"time"

	"github.com/google/uuid"
)

// Borrower represents a domain entity for a loan borrower
type Borrower struct {
	ID          string
	FullName    string
	Email       string
	PhoneNumber string
	IDNumber    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBorrower(fullName, email, phoneNumber, idNumber string) *Borrower {
	now := time.Now()
	return &Borrower{
		ID:          uuid.New().String(),
		FullName:    fullName,
		Email:       email,
		PhoneNumber: phoneNumber,
		IDNumber:    idNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
