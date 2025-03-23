package migrations_models

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

type Loan struct {
	ID                  string   `gorm:"type:uuid;primary_key"`
	BorrowerID          string   `gorm:"type:uuid;index:idx_loan_borrower_id;not null"`
	Borrower            Borrower `gorm:"foreignKey:BorrowerID"`
	Amount              float64  `gorm:"type:decimal(20,2);not null"`
	Rate                float64  `gorm:"type:decimal(5,2);not null"` // Total interest rate for borrower
	ROI                 float64  `gorm:"type:decimal(5,2);not null"` // Return on Investment for investors
	Status              string   `gorm:"index:idx_loan_status;type:varchar(20);not null"`
	SurveyDocumentID    string   `gorm:"type:uuid;index:idx_survey_loan_document_id"`
	ApprovalDate        *time.Time
	ApprovedBy          string `gorm:"type:uuid;index;default:null"`
	InvestmentDate      *time.Time
	DisbursementDate    *time.Time
	DisbursedBy         string       `gorm:"type:uuid;index;default:null"`
	AgreementDocumentID string       `gorm:"type:uuid;index:idx_agreement_loan_document_id"`
	StatusTransitions   JSON         `gorm:"type:jsonb"`
	LoanLenders         []LoanLender `gorm:"foreignKey:LoanID"`
	CreatedAt           time.Time    `gorm:"index"`
	UpdatedAt           time.Time
}

type JSON []loan.StatusTransition
