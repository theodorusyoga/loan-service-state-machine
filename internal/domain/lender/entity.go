package lender

import (
	"time"

	"github.com/google/uuid"
)

// Lender represents a domain entity for a loan lender/investor
type Lender struct {
	ID          string
	FullName    string
	Email       string
	PhoneNumber string
	IDNumber    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewLender(fullName, email, phoneNumber, idNumber string) *Lender {
	now := time.Now()
	return &Lender{
		ID:          uuid.New().String(),
		FullName:    fullName,
		Email:       email,
		PhoneNumber: phoneNumber,
		IDNumber:    idNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
