package model

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

func (LoanModel) TableName() string {
	return "loans"
}

type LoanModel struct {
	ID                string `gorm:"type:uuid;primary_key"`
	Amount            float64
	Status            string `gorm:"index;type:varchar(20)"`
	ApprovalDate      *time.Time
	ApprovedBy        string
	InvestmentDate    *time.Time
	DisbursementDate  *time.Time
	DisbursedBy       string
	StatusTransitions JSON      `gorm:"type:jsonb"` // Store as JSONB for CockroachDB
	CreatedAt         time.Time `gorm:"index"`
	UpdatedAt         time.Time
}

type JSON []loan.StatusTransition

func (m *LoanModel) ToEntity() *loan.Loan {
	return &loan.Loan{
		ID:                m.ID,
		Amount:            m.Amount,
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

func FromEntity(l *loan.Loan) *LoanModel {
	return &LoanModel{
		ID:                l.ID,
		Amount:            l.Amount,
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

func (m *LoanModel) ToDomain() *loan.Loan {
	return &loan.Loan{
		ID:                m.ID,
		Amount:            m.Amount,
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
