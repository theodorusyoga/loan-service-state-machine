package employee

import (
	"context"
)

// Repository defines the data access interface for employees
type Repository interface {
	Get(ctx context.Context, id string) (*Employee, error)
	Save(ctx context.Context, employee *Employee) error
	Create(ctx context.Context, Employee *Employee) error
	// TODO: Implement the following methods
	List(ctx context.Context, filter EmployeeFilter) ([]*Employee, error)
	// Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter EmployeeFilter) (int64, error)
}

type EmployeeFilter struct {
	FullName    *string
	Email       *string
	PhoneNumber *string
	IDNumber    *string
	Page        int
	PageSize    int
}

func (f *EmployeeFilter) WithDefaults() *EmployeeFilter {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
	return f
}
