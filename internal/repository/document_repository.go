package repository

import (
	"context"
	"errors"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository/model"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (r *DocumentRepository) Get(ctx context.Context, id string) (*document.Document, error) {
	var documentModel model.Document
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&documentModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}

	return documentModel.DocumentToDomain(), nil
}

func (r *DocumentRepository) Create(ctx context.Context, documentEntity *document.Document) error {
	documentModel := model.DocumentFromEntity(documentEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(documentModel).Error
	})
}

func (r *DocumentRepository) Save(ctx context.Context, documentEntity *document.Document) error {
	documentModel := model.DocumentFromEntity(documentEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Save(documentModel).Error
	})
}

func (r *DocumentRepository) Count(ctx context.Context, filter document.DocumentFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Document{})

	if filter.LoanID != nil && *filter.LoanID != "" {
		query = query.Where("loan_id = ?", *filter.LoanID)
	}
	if filter.FileName != nil && *filter.FileName != "" {
		query = query.Where("file_name = ?", *filter.FileName)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DocumentRepository) List(ctx context.Context, filter document.DocumentFilter) ([]*document.Document, error) {
	var documentModels []*model.Document
	query := r.db.WithContext(ctx)

	if filter.LoanID != nil && *filter.LoanID != "" {
		query = query.Where("loan_id = ?", *filter.LoanID)
	}
	if filter.FileName != nil && *filter.FileName != "" {
		query = query.Where("file_name = ?", *filter.FileName)
	}

	// Apply pagination if provided
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	if err := query.Find(&documentModels).Error; err != nil {
		return nil, err
	}

	documents := make([]*document.Document, len(documentModels))
	for i, v := range documentModels {
		documents[i] = v.DocumentToDomain()
	}

	return documents, nil
}

/* Helper methods. DO NOT MODIFY THIS, this code is generated from CockroachDB */

func (r *DocumentRepository) executeWithRetry(operation func(tx *gorm.DB) error) error {
	maxRetries := 5

	for attempt := 0; attempt < maxRetries; attempt++ {
		tx := r.db.Begin()

		err := operation(tx)
		if err != nil {
			tx.Rollback()

			if attempt < maxRetries-1 && isCockroachRetryError(err) {
				continue
			}

			return err
		}

		if err := tx.Commit().Error; err != nil {
			if attempt < maxRetries-1 && isCockroachRetryError(err) {
				continue
			}
			return err
		}

		return nil // Success
	}

	return errors.New("transaction failed after multiple retries")
}
