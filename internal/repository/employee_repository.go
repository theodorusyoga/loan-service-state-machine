package repository

import (
	"context"
	"errors"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository/model"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) Get(ctx context.Context, id string) (*employee.Employee, error) {
	var employeeModel model.Employee
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&employeeModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}

	return employeeModel.EmployeeToDomain(), nil
}

func (r *EmployeeRepository) Create(ctx context.Context, employeeEntity *employee.Employee) error {
	employeeModel := model.EmployeeFromEntity(employeeEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(employeeModel).Error
	})
}

func (r *EmployeeRepository) Save(ctx context.Context, employeeEntity *employee.Employee) error {
	employeeModel := model.EmployeeFromEntity(employeeEntity)

	// Use CockroachDB transaction retry logic
	return r.executeWithRetry(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Save(employeeModel).Error
	})
}

func (r *EmployeeRepository) List(ctx context.Context, filter employee.EmployeeFilter) ([]*employee.Employee, error) {
	var employeeModels []*model.Employee
	query := r.db.WithContext(ctx)

	if filter.FullName != nil && *filter.FullName != "" {
		query = query.Where("full_name = ?", *filter.FullName)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", filter.PhoneNumber)
	}

	if filter.IDNumber != nil && *filter.IDNumber != "" {
		query = query.Where("id_number = ?", filter.IDNumber)
	}

	if err := query.Find(&employeeModels).Error; err != nil {
		return nil, err
	}

	employees := make([]*employee.Employee, len(employeeModels))
	for i, v := range employeeModels {
		employees[i] = v.EmployeeToDomain()
	}

	return employees, nil
}

/* Helper methods. DO NOT MODIFY THIS, this code is generated from CockroachDB */

func (r *EmployeeRepository) executeWithRetry(operation func(tx *gorm.DB) error) error {
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
