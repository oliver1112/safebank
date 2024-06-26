package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("email duplicate")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (ud *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}

func (ud *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}

func (ud *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now

	err := ud.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErr uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErr {
			// email conflict
			return ErrUserDuplicateEmail
		}
	}

	return err
}

// User is the mapping to database form, also called entity, model, PO(persistent object)
type User struct {
	ID      int64  `gorm:"primaryKey,autoIncrement" json:"id"`
	Email   string `gorm:"unique" json:"email"`
	FName   string `json:"f_name"`
	LName   string `json:"l_name"`
	Country string `json:"country"`

	State  string `json:"state"`
	City   string `json:"city"`
	Street string `json:"street"`
	Apart  string `json:"party"`
	Zip    string `json:"zip"`

	Password string
	Ctime    int64 `json:"ctime"`
	Utime    int64 `json:"utime"`

	Account []Account
}

func (User) TableName() string {
	return "wsj_user"
}
