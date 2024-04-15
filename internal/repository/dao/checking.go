package dao

import (
	"context"
	"gorm.io/gorm"
)

type CheckingDAO struct {
	db *gorm.DB
}

type Checking struct {
	AccountID     int64 `gorm:"primaryKey"`
	Account       Account
	ServiceCharge float64
	Amount        int64
}

func (sd *CheckingDAO) Insert(ctx context.Context, c Checking) error {
	err := sd.db.WithContext(ctx).Create(&c).Error
	return err
}
