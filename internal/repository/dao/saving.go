package dao

import (
	"context"
	"gorm.io/gorm"
)

type SavingDAO struct {
	db *gorm.DB
}

type Saving struct {
	AccountID    int64 `gorm:"primaryKey"`
	Account      Account
	InterestRate float64
	Amount       float64
}

func NewSavingDao(db *gorm.DB) *SavingDAO {
	return &SavingDAO{
		db: db,
	}
}

func (sd *SavingDAO) Insert(ctx context.Context, s Saving) error {
	err := sd.db.WithContext(ctx).Create(&s).Error
	return err
}
