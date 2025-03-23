package migrations_models

import (
	"time"
)

type Lender struct {
	ID          string       `gorm:"type:uuid;primary_key"`
	FullName    string       `gorm:"type:varchar(100);not null"`
	Email       string       `gorm:"type:varchar(100);uniqueIndex:uni_lenders_email;not null"`
	PhoneNumber string       `gorm:"type:varchar(20);not null"`
	IDNumber    string       `gorm:"type:varchar(50);uniqueIndex:uni_lenders_id_number;not null"`
	LoanLenders []LoanLender `gorm:"foreignKey:LenderID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
