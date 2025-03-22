package migrations_models

import (
	"time"
)

type Borrower struct {
	ID          string `gorm:"type:uuid;primary_key"`
	FullName    string `gorm:"type:varchar(100);not null"`
	Email       string `gorm:"type:varchar(100);uniqueIndex:uni_borrowers_email;not null"`
	PhoneNumber string `gorm:"type:varchar(20);not null"`
	IDNumber    string `gorm:"type:varchar(50);uniqueIndex:uni_borrowers_id_number;not null"`
	Loans       []Loan `gorm:"foreignKey:BorrowerID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
