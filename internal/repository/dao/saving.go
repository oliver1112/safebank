package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SavingDAO struct {
	Db *gorm.DB
}

type Saving struct {
	AccountID    int64   `gorm:"primaryKey" json:"account_id"`
	Account      Account `json:"account"`
	InterestRate float64 `json:"interest_rate"`
	Amount       float64 `json:"amount"`
}

func NewSavingDao(db *gorm.DB) *SavingDAO {
	return &SavingDAO{
		Db: db,
	}
}

func (sd *SavingDAO) GetSaving(ctx *gin.Context, accountId int64) (Saving, error) {
	var saving Saving
	err := sd.Db.WithContext(ctx).Where("account_id = ?", accountId).First(&saving).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return saving, err
}

func (sd *SavingDAO) CreateOrUpdate(ctx *gin.Context, data Saving) (Saving, error) {
	where := Saving{
		AccountID: data.AccountID,
	}
	var saving Saving
	err := sd.Db.Where(where).Assign(data).FirstOrCreate(&saving).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return saving, err
}

func (sd *SavingDAO) Insert(ctx context.Context, s Saving) error {
	err := sd.Db.WithContext(ctx).Create(&s).Error
	return err
}
