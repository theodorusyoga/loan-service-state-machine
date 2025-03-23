package borrower

import (
	"time"

	"github.com/google/uuid"
)

// Borrower represents a domain entity for a loan borrower
type Borrower struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IDNumber    string    `json:"id_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
