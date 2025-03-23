package migrations_models

import (
	"time"
)

type LoanLender struct {
	ID         string    `gorm:"type:uuid;primary_key"`
	LoanID     string    `gorm:"type:uuid;index:idx_loan_lender_loan_id;not null"`
	Loan       Loan      `gorm:"foreignKey:LoanID"`
	LenderID   string    `gorm:"type:uuid;index:idx_loan_lender_lender_id;not null"`
	Lender     Lender    `gorm:"foreignKey:LenderID"`
	Amount     float64   `gorm:"type:decimal(20,2);not null"` // Amount invested by this lender
	InvestedAt time.Time `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
