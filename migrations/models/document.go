package migrations_models

import (
	"time"
)

type Document struct {
	ID        string `gorm:"type:uuid;primary_key"`
	FileName  string `gorm:"type:varchar(100)"` // use dummy filename for now
	CreatedAt time.Time
	UpdatedAt time.Time

	SurveyLoans    []Loan `gorm:"foreignKey:SurveyDocumentID"`
	AgreementLoans []Loan `gorm:"foreignKey:AgreementDocumentID"`
}
