package dao

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (sd *SavingDAO) GetSaving(ctx *gin.Context, userId int64) (Saving, error) {
	var saving Saving
	err := sd.db.WithContext(ctx).Where("user_id = ?", userId).First(&saving).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return saving, err
}

func (sd *SavingDAO) CreateOrUpdate(ctx context.Context, data Saving) (Saving, error) {
	where := Saving{
		AccountID: data.AccountID,
	}
	var saving Saving
	err := sd.db.Where(where).Assign(data).FirstOrCreate(&saving).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return saving, err
}

func (sd *SavingDAO) Insert(ctx context.Context, s Saving) error {
	err := sd.db.WithContext(ctx).Create(&s).Error
	return err
}
