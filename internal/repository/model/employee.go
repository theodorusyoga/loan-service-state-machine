package model

import (
	"time"

	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
)

func (Employee) TableName() string {
	return "employees"
}

type Employee struct {
	ID          string `gorm:"type:uuid;primary_key"`
	FullName    string `gorm:"type:varchar(100)"`
	Email       string `gorm:"type:varchar(100);uniqueIndex"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	IDNumber    string `gorm:"type:varchar(50);uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Employee) EmployeeToEntity() *employee.Employee {
	return &employee.Employee{
		ID:          m.ID,
		FullName:    m.FullName,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		IDNumber:    m.IDNumber,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func EmployeeFromEntity(e *employee.Employee) *Employee {
	return &Employee{
		ID:          e.ID,
		FullName:    e.FullName,
		Email:       e.Email,
		PhoneNumber: e.PhoneNumber,
		IDNumber:    e.IDNumber,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (m *Employee) EmployeeToDomain() *employee.Employee {
	return &employee.Employee{
		ID:          m.ID,
		FullName:    m.FullName,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		IDNumber:    m.IDNumber,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
