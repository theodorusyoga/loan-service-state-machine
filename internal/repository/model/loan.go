package model

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

func (Loan) TableName() string {
	return "loans"
}

type Loan struct {
	ID                string `gorm:"type:uuid;primary_key"`
	BorrowerID        string `gorm:"type:uuid;index"`
	Amount            float64
	Rate              float64
	ROI               float64
	Status            string `gorm:"index;type:varchar(20)"`
	ApprovalDate      *time.Time
	ApprovedBy        *string
	InvestmentDate    *time.Time
	DisbursementDate  *time.Time
	DisbursedBy       *string
	StatusTransitions JSON      `gorm:"type:jsonb"` // Store as JSONB for CockroachDB
	CreatedAt         time.Time `gorm:"index"`
	UpdatedAt         time.Time
}

type JSON []loan.StatusTransition

func (m *Loan) LoanToEntity() *loan.Loan {
	return &loan.Loan{
		ID:                m.ID,
		BorrowerID:        m.BorrowerID,
		Amount:            m.Amount,
		Rate:              m.Rate,
		ROI:               m.ROI,
		Status:            loan.Status(m.Status),
		ApprovalDate:      m.ApprovalDate,
		ApprovedBy:        m.ApprovedBy,
		InvestmentDate:    m.InvestmentDate,
		DisbursementDate:  m.DisbursementDate,
		DisbursedBy:       m.DisbursedBy,
		StatusTransitions: []loan.StatusTransition(m.StatusTransitions),
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func LoanFromEntity(l *loan.Loan) *Loan {
	return &Loan{
		ID:                l.ID,
		BorrowerID:        l.BorrowerID,
		Amount:            l.Amount,
		Rate:              l.Rate,
		ROI:               l.ROI,
		Status:            string(l.Status),
		ApprovalDate:      l.ApprovalDate,
		ApprovedBy:        l.ApprovedBy,
		InvestmentDate:    l.InvestmentDate,
		DisbursementDate:  l.DisbursementDate,
		DisbursedBy:       l.DisbursedBy,
		StatusTransitions: JSON(l.StatusTransitions),
		CreatedAt:         l.CreatedAt,
		UpdatedAt:         l.UpdatedAt,
	}
}

func (m *Loan) LoanToDomain() *loan.Loan {
	return &loan.Loan{
		ID:                m.ID,
		BorrowerID:        m.BorrowerID,
		Amount:            m.Amount,
		Rate:              m.Rate,
		ROI:               m.ROI,
		Status:            loan.Status(m.Status),
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		ApprovalDate:      m.ApprovalDate,
		ApprovedBy:        m.ApprovedBy,
		InvestmentDate:    m.InvestmentDate,
		DisbursementDate:  m.DisbursementDate,
		StatusTransitions: m.StatusTransitions,
	}
}
