package document

import (
	"time"

	"github.com/google/uuid"
)

// Document represents a domain entity for a loan document (e.g. agreement letter)
type Document struct {
	ID        string
	FileName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDocument(loanID, fileName string) *Document {
	now := time.Now()
	return &Document{
		ID:        uuid.New().String(),
		FileName:  fileName,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
