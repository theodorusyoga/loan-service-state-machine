package repository

import (
	"context"
	"errors"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/borrower"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository/model"
	"gorm.io/gorm"
)

type BorrowerRepository struct {
	db *gorm.DB
}

func NewBorrowerRepository(db *gorm.DB) *BorrowerRepository {
	return &BorrowerRepository{
		db: db,
	}
}

func (r *BorrowerRepository) Get(ctx context.Context, id string) (*borrower.Borrower, error) {
	var borrowerModel model.Borrower
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&borrowerModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrower not found")
		}
		return nil, err
	}

	return borrowerModel.BorrowerToDomain(), nil
}

func (r *BorrowerRepository) Create(ctx context.Context, borrowerEntity *borrower.Borrower) error {
	borrowerModel := model.BorrowerFromEntity(borrowerEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(borrowerModel).Error
	})
}

func (r *BorrowerRepository) Save(ctx context.Context, borrowerEntity *borrower.Borrower) error {
	borrowerModel := model.BorrowerFromEntity(borrowerEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Save(borrowerModel).Error
	})
}

/* Helper methods. DO NOT MODIFY THIS, this code is generated from CockroachDB */

func (r *BorrowerRepository) executeWithRetry(operation func(tx *gorm.DB) error) error {
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
