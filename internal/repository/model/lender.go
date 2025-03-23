package model

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
)

func (Lender) TableName() string {
	return "lenders"
}

type Lender struct {
	ID          string `gorm:"type:uuid;primary_key"`
	FullName    string `gorm:"type:varchar(100)"`
	Email       string `gorm:"type:varchar(100);uniqueIndex"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	IDNumber    string `gorm:"type:varchar(50);uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Lender) LenderToEntity() *lender.Lender {
	return &lender.Lender{
		ID:          m.ID,
		FullName:    m.FullName,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		IDNumber:    m.IDNumber,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func LenderFromEntity(l *lender.Lender) *Lender {
	return &Lender{
		ID:          l.ID,
		FullName:    l.FullName,
		Email:       l.Email,
		PhoneNumber: l.PhoneNumber,
		IDNumber:    l.IDNumber,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

func (m *Lender) LenderToDomain() *lender.Lender {
	return &lender.Lender{
		ID:          m.ID,
		FullName:    m.FullName,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		IDNumber:    m.IDNumber,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
