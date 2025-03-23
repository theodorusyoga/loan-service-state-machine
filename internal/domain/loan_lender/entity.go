package loanlender

import (
	"time"

	"github.com/google/uuid"
)

// LoanLender represents a domain entity for the relationship between a loan and a lender
type LoanLender struct {
	ID         string
	LoanID     string
	LenderID   string
	Amount     float64
	InvestedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewLoanLender(loanID, lenderID string, amount float64) *LoanLender {
	now := time.Now()
	return &LoanLender{
		ID:         uuid.New().String(),
		LoanID:     loanID,
		LenderID:   lenderID,
		Amount:     amount,
		InvestedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
