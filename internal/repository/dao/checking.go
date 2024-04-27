package dao

import (
	"context"
	"github.com/gin-gonic/gin"
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

func NewCheckingDao(db *gorm.DB) *CheckingDAO {
	return &CheckingDAO{
		db: db,
	}
}

func (cd *CheckingDAO) Insert(ctx context.Context, c Checking) error {
	err := cd.db.WithContext(ctx).Create(&c).Error
	return err
}

func (cd *CheckingDAO) GetChecking(ctx *gin.Context, accountId int64) (Checking, error) {
	var checking Checking
	err := cd.db.WithContext(ctx).Where("account_id = ?", accountId).First(&checking).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return checking, err
}

func (cd *CheckingDAO) CreateOrUpdate(ctx *gin.Context, data Checking) (Checking, error) {
	where := Checking{
		AccountID: data.AccountID,
	}
	var checking Checking
	err := cd.db.Where(where).Assign(data).FirstOrCreate(&checking).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return checking, err
}
