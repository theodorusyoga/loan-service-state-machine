package migrations_models

import (
	"time"
)

type Document struct {
	ID        string `gorm:"type:uuid;primary_key"`
	LoanID    string `gorm:"type:uuid;uniqueIndex:uni_documents_loan_id"`
	Loan      Loan   `gorm:"foreignKey:LoanID"`
	FileName  string `gorm:"type:varchar(100)"` // use dummy filename for now
	CreatedAt time.Time
	UpdatedAt time.Time
}
