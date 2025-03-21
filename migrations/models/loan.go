package migrations_models

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

type Loan struct {
	ID                string `gorm:"type:uuid;primary_key"`
	Amount            float64
	Status            string `gorm:"index;type:varchar(20)"`
	ApprovalDate      *time.Time
	ApprovedBy        string
	InvestmentDate    *time.Time
	DisbursementDate  *time.Time
	DisbursedBy       string
	StatusTransitions JSON      `gorm:"type:jsonb"`
	CreatedAt         time.Time `gorm:"index"`
	UpdatedAt         time.Time
}

type JSON []loan.StatusTransition
