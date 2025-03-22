package model

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/borrower"
)

func (Borrower) TableName() string {
	return "borrowers"
}

type Borrower struct {
	ID          string    `gorm:"type:uuid;primary_key"`
	FullName    string    `gorm:"type:varchar(100)"`
	Email       string    `gorm:"type:varchar(100);index"`
	PhoneNumber string    `gorm:"type:varchar(20)"`
	IDNumber    string    `gorm:"type:varchar(50);index"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
}

func (m *Borrower) BorrowerToEntity() *borrower.Borrower {
	return &borrower.Borrower{
		ID:          m.ID,
		FullName:    m.FullName,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		IDNumber:    m.IDNumber,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func BorrowerFromEntity(b *borrower.Borrower) *Borrower {
	return &Borrower{
		ID:          b.ID,
		FullName:    b.FullName,
		Email:       b.Email,
		PhoneNumber: b.PhoneNumber,
		IDNumber:    b.IDNumber,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

// ToDomain maintains consistency with the loan.go file by providing an alias for ToEntity
func (m *Borrower) BorrowerToDomain() *borrower.Borrower {
	return m.BorrowerToEntity()
}
