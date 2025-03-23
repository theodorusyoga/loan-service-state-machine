package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository/model"
	"gorm.io/gorm"
)

type LoanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{
		db: db,
	}
}

func (r *LoanRepository) Get(ctx context.Context, id string) (*loan.Loan, error) {
	var loanModel model.Loan
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&loanModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}

	return loanModel.LoanToDomain(), nil
}

func (r *LoanRepository) Create(ctx context.Context, loanEntity *loan.Loan) error {
	loanModel := model.LoanFromEntity(loanEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(loanModel).Error
	})
}

func (r *LoanRepository) Save(ctx context.Context, loanEntity *loan.Loan) error {
	loanModel := model.LoanFromEntity(loanEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Save(loanModel).Error
	})
}

func (r *LoanRepository) Count(ctx context.Context, filter loan.LoanFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Loan{})

	if filter.MaxAmount != nil && *filter.MaxAmount != 0 {
		query = query.Where("amount < ?", *filter.MaxAmount)
	}

	if filter.MinAmount != nil && *filter.MinAmount != 0 {
		query = query.Where("amount > ?", *filter.MinAmount)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *LoanRepository) List(ctx context.Context, filter loan.LoanFilter) ([]*loan.Loan, error) {
	var loanModels []*model.Loan
	query := r.db.WithContext(ctx).Model(&model.Loan{})

	if filter.MaxAmount != nil && *filter.MaxAmount != 0 {
		query = query.Where("amount < ?", *filter.MaxAmount)
	}

	if filter.MinAmount != nil && *filter.MinAmount != 0 {
		query = query.Where("amount > ?", *filter.MinAmount)
	}

	if err := query.Find(&loanModels).Error; err != nil {
		return nil, err
	}

	var loans []*loan.Loan
	for _, loanModel := range loanModels {
		loans = append(loans, loanModel.LoanToDomain())
	}

	return loans, nil
}

/* Helper methods. DO NOT MODIFY THIS, this code is generated from CockroachDB */

func (r *LoanRepository) executeWithRetry(operation func(tx *gorm.DB) error) error {
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

func isCockroachRetryError(err error) bool {
	// CockroachDB retry error codes typically contain 40001 or
	// message about transaction retry
	return err != nil && (errors.Is(err, gorm.ErrInvalidTransaction) ||
		containsAny(err.Error(), []string{"40001", "retry transaction", "restart transaction"}))
}

func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
