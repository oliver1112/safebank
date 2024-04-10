package dao

import "gorm.io/gorm"

type CheckingDAO struct {
	db *gorm.DB
}

type Checking struct {
	AccountID     int64 `gorm:"primaryKey"`
	Account       Account
	ServiceCharge float64
	Amount        int64
}
