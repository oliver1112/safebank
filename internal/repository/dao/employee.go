package dao

import (
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// User is the mapping to database form, also called entity, model, PO(persistent object)
type Employee struct {
	ID      int64  `gorm:"primaryKey,autoIncrement"`
	Email   string `gorm:"unique"`
	FName   string
	LName   string
	Country string

	State  string
	City   string
	Street string
	Apart  string
	Zip    string

	Password string
	Ctime    int64
	Utime    int64
}

type EmployeeDAO struct {
	db *gorm.DB
}

func NewEmployeeDao(db *gorm.DB) *EmployeeDAO {
	return &EmployeeDAO{
		db: db,
	}
}

func (e *EmployeeDAO) CreateOrUpdate(data Employee) (Employee, error) {
	where := Employee{
		Email: cast.ToString(data.Email),
	}
	var employee Employee
	err := e.db.Where(where).Assign(data).FirstOrCreate(&employee).Error
	return employee, err
}

func (e *EmployeeDAO) FindByEmail(ctx context.Context, email string) (Employee, error) {
	var employee Employee
	err := e.db.WithContext(ctx).Where("email = ?", email).First(&employee).Error
	return employee, err
}
