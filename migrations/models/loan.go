package migrations_models

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

type Loan struct {
	ID                string   `gorm:"type:uuid;primary_key"`
	BorrowerID        string   `gorm:"type:uuid;index:idx_loan_borrower_id;not null"`
	Borrower          Borrower `gorm:"foreignKey:BorrowerID"`
	Amount            float64  `gorm:"type:decimal(20,2);not null"`
	Rate              float64  `gorm:"type:decimal(5,2);not null"` // Total interest rate for borrower
	ROI               float64  `gorm:"type:decimal(5,2);not null"` // Return on Investment for investors
	Status            string   `gorm:"index:idx_loan_status;type:varchar(20);not null"`
	ApprovalDate      *time.Time
	ApprovedBy        string
	InvestmentDate    *time.Time
	DisbursementDate  *time.Time
	DisbursedBy       string
	AgreementLetterID string    `gorm:"type:uuid;index:idx_loan_agreement_letter_id"`
	StatusTransitions JSON      `gorm:"type:jsonb"`
	CreatedAt         time.Time `gorm:"index"`
	UpdatedAt         time.Time
}

type JSON []loan.StatusTransition
