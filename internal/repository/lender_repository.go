package repository

import (
	"context"
	"errors"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository/model"
	"gorm.io/gorm"
)

type LenderRepository struct {
	db *gorm.DB
}

func NewLenderRepository(db *gorm.DB) *LenderRepository {
	return &LenderRepository{
		db: db,
	}
}

func (r *LenderRepository) Get(ctx context.Context, id string) (*lender.Lender, error) {
	var lenderModel model.Lender
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&lenderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lender not found")
		}
		return nil, err
	}

	return lenderModel.LenderToDomain(), nil
}

func (r *LenderRepository) Create(ctx context.Context, lenderEntity *lender.Lender) error {
	lenderModel := model.LenderFromEntity(lenderEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(lenderModel).Error
	})
}

func (r *LenderRepository) Save(ctx context.Context, lenderEntity *lender.Lender) error {
	lenderModel := model.LenderFromEntity(lenderEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Save(lenderModel).Error
	})
}

func (r *LenderRepository) Count(ctx context.Context, filter lender.LenderFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Lender{})

	if filter.FullName != nil && *filter.FullName != "" {
		query = query.Where("full_name = ?", *filter.FullName)
	}
	if filter.Email != nil && *filter.Email != "" {
		query = query.Where("email = ?", *filter.Email)
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", *filter.PhoneNumber)
	}
	if filter.IDNumber != nil && *filter.IDNumber != "" {
		query = query.Where("id_number = ?", *filter.IDNumber)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *LenderRepository) List(ctx context.Context, filter lender.LenderFilter) ([]*lender.Lender, error) {
	var lenderModels []*model.Lender

	query := r.db.WithContext(ctx)

	if filter.FullName != nil && *filter.FullName != "" {
		query = query.Where("full_name = ?", *filter.FullName)
	}
	if filter.Email != nil && *filter.Email != "" {
		query = query.Where("email = ?", *filter.Email)
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", *filter.PhoneNumber)
	}
	if filter.IDNumber != nil && *filter.IDNumber != "" {
		query = query.Where("id_number = ?", *filter.IDNumber)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		query = query.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	}

	if err := query.Find(&lenderModels).Error; err != nil {
		return nil, err
	}

	var lenders []*lender.Lender
	for _, lenderModel := range lenderModels {
		lenders = append(lenders, lenderModel.LenderToDomain())
	}

	return lenders, nil
}

/* Helper methods. DO NOT MODIFY THIS, this code is generated from CockroachDB */

func (r *LenderRepository) executeWithRetry(operation func(tx *gorm.DB) error) error {
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
