package dao

import (
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// User is the mapping to database form, also called entity, model, PO(persistent object)
type Employee struct {
	ID      int64  `gorm:"primaryKey,autoIncrement" json:"id"`
	Email   string `gorm:"unique" json:"email"`
	FName   string `json:"fname"`
	LName   string `json:"lname"`
	Country string `json:"country"`

	State  string `json:"state"`
	City   string `json:"city"`
	Street string `json:"street"`
	Apart  string `json:"party"`
	Zip    string `json:"zip"`

	Password string `json:"password"`
	Ctime    int64  `json:"ctime"`
	Utime    int64  `json:"utime"`
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
