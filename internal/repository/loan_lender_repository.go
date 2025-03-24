package repository

import (
	"context"
	"errors"

	loanlender "github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan_lender"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository/model"
	"gorm.io/gorm"
)

type LoanLenderRepository struct {
	db *gorm.DB
}

func NewLoanLenderRepository(db *gorm.DB) *LoanLenderRepository {
	return &LoanLenderRepository{
		db: db,
	}
}

func (r *LoanLenderRepository) Get(ctx context.Context, id string) (*loanlender.LoanLender, error) {
	var loanLenderModel model.LoanLender
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&loanLenderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan-lender relationship not found")
		}
		return nil, err
	}

	return loanLenderModel.LoanLenderToDomain(), nil
}

func (r *LoanLenderRepository) GetByLoanID(ctx context.Context, loanID string) ([]*loanlender.LoanLender, error) {
	var loanLenderModels []*model.LoanLender
	if err := r.db.WithContext(ctx).Where("loan_id = ?", loanID).Find(&loanLenderModels).Error; err != nil {
		return nil, err
	}

	var loanLenders []*loanlender.LoanLender
	for _, loanLenderModel := range loanLenderModels {
		loanLenders = append(loanLenders, loanLenderModel.LoanLenderToDomain())
	}

	return loanLenders, nil
}

func (r *LoanLenderRepository) GetByLenderID(ctx context.Context, lenderID string) ([]*loanlender.LoanLender, error) {
	var loanLenderModels []*model.LoanLender
	if err := r.db.WithContext(ctx).Where("lender_id = ?", lenderID).Find(&loanLenderModels).Error; err != nil {
		return nil, err
	}

	var loanLenders []*loanlender.LoanLender
	for _, loanLenderModel := range loanLenderModels {
		loanLenders = append(loanLenders, loanLenderModel.LoanLenderToDomain())
	}

	return loanLenders, nil
}

func (r *LoanLenderRepository) Create(ctx context.Context, loanLenderEntity *loanlender.LoanLender) error {
	loanLenderModel := model.LoanLenderFromEntity(loanLenderEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(loanLenderModel).Error
	})
}

func (r *LoanLenderRepository) Save(ctx context.Context, loanLenderEntity *loanlender.LoanLender) error {
	loanLenderModel := model.LoanLenderFromEntity(loanLenderEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Save(loanLenderModel).Error
	})
}

func (r *LoanLenderRepository) Count(ctx context.Context, filter loanlender.LoanLenderFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.LoanLender{})

	query = r.applyFilter(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *LoanLenderRepository) List(ctx context.Context, filter loanlender.LoanLenderFilter) ([]*loanlender.LoanLender, error) {
	var loanLenderModels []*model.LoanLender

	query := r.db.WithContext(ctx)
	query = r.applyFilter(query, filter)

	// Apply pagination
	filter.WithDefaults()
	query = query.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)

	if err := query.Find(&loanLenderModels).Error; err != nil {
		return nil, err
	}

	var loanLenders []*loanlender.LoanLender
	for _, loanLenderModel := range loanLenderModels {
		loanLenders = append(loanLenders, loanLenderModel.LoanLenderToDomain())
	}

	return loanLenders, nil
}

func (r *LoanLenderRepository) applyFilter(query *gorm.DB, filter loanlender.LoanLenderFilter) *gorm.DB {
	if filter.LoanID != nil && *filter.LoanID != "" {
		query = query.Where("loan_id = ?", *filter.LoanID)
	}
	if filter.LenderID != nil && *filter.LenderID != "" {
		query = query.Where("lender_id = ?", *filter.LenderID)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}
	if filter.InvestedFrom != nil {
		query = query.Where("invested_at >= ?", *filter.InvestedFrom)
	}
	if filter.InvestedTo != nil {
		query = query.Where("invested_at <= ?", *filter.InvestedTo)
	}

	return query
}

/* Helper methods. DO NOT MODIFY THIS, this code is generated from CockroachDB */

func (r *LoanLenderRepository) executeWithRetry(operation func(tx *gorm.DB) error) error {
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
