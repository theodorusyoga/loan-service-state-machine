package model

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
)

func (Document) TableName() string {
	return "documents"
}

type Document struct {
	ID        string `gorm:"type:uuid;primary_key"`
	FileName  string `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Document) DocumentToEntity() *document.Document {
	return &document.Document{
		ID:        m.ID,
		FileName:  m.FileName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func DocumentFromEntity(d *document.Document) *Document {
	return &Document{
		ID:        d.ID,
		FileName:  d.FileName,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (m *Document) DocumentToDomain() *document.Document {
	return &document.Document{
		ID:        m.ID,
		FileName:  m.FileName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
