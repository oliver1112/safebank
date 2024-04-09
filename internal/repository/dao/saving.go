package dao

import "gorm.io/gorm"

type SavingDAO struct {
	db *gorm.DB
}

type Saving struct {
	AccountID    int `gorm:"primaryKey"`
	Account      Account
	InterestRate float64
}
