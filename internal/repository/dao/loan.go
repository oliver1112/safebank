package dao

import "gorm.io/gorm"

type LoanDAO struct {
	db *gorm.DB
}

type Loan struct {
	AccountID int `gorm:"primaryKey"`
	Account   Account
	Rate      float64
	Amount    float64
	Month     int
	Payment   float64
	Type      string
}
