package migrations_models

import (
	"time"
)

type Employee struct {
	ID             string `gorm:"type:uuid;primary_key"`
	FullName       string `gorm:"type:varchar(100);not null"`
	Email          string `gorm:"type:varchar(100);uniqueIndex:uni_employees_email;not null"`
	PhoneNumber    string `gorm:"type:varchar(20);not null"`
	IDNumber       string `gorm:"type:varchar(50);uniqueIndex:uni_employees_id_number;not null"`
	ApprovedLoans  []Loan `gorm:"foreignKey:ApprovedBy"`
	DisbursedLoans []Loan `gorm:"foreignKey:DisbursedBy"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
