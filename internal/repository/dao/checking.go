package dao

import "gorm.io/gorm"

type CheckingDAO struct {
	db *gorm.DB
}

type Checking struct {
	AccountID     int `gorm:"primaryKey"`
	Account       Account
	ServiceCharge float64
}
