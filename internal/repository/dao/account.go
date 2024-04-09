package dao

import "gorm.io/gorm"

type AccountDAO struct {
	db *gorm.DB
}

type Account struct {
	ID   int64 `gorm:"primaryKey,autoIncrement"`
	name string

	Street      string
	City        string
	State       string
	Zip         string
	AccountType string
	UserID      int64

	Ctime int64
	Utime int64
}
