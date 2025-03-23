package model

import (
	"time"

	loanlender "github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan_lender"
)

func (LoanLender) TableName() string {
	return "loan_lenders"
}

type LoanLender struct {
	ID         string `gorm:"type:uuid;primary_key"`
	LoanID     string `gorm:"type:uuid;index"`
	LenderID   string `gorm:"type:uuid;index"`
	Amount     float64
	InvestedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *LoanLender) LoanLenderToEntity() *loanlender.LoanLender {
	return &loanlender.LoanLender{
		ID:         m.ID,
		LoanID:     m.LoanID,
		LenderID:   m.LenderID,
		Amount:     m.Amount,
		InvestedAt: m.InvestedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func LoanLenderFromEntity(ll *loanlender.LoanLender) *LoanLender {
	return &LoanLender{
		ID:         ll.ID,
		LoanID:     ll.LoanID,
		LenderID:   ll.LenderID,
		Amount:     ll.Amount,
		InvestedAt: ll.InvestedAt,
		CreatedAt:  ll.CreatedAt,
		UpdatedAt:  ll.UpdatedAt,
	}
}

func (m *LoanLender) LoanLenderToDomain() *loanlender.LoanLender {
	return &loanlender.LoanLender{
		ID:         m.ID,
		LoanID:     m.LoanID,
		LenderID:   m.LenderID,
		Amount:     m.Amount,
		InvestedAt: m.InvestedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
