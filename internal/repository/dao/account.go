package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type AccountDAO struct {
	db *gorm.DB
}

type Account struct {
	ID          int64 `gorm:"primaryKey,autoIncrement"`
	Name        string
	Street      string
	City        string
	State       string
	Zip         string
	Apart       string
	AccountType string
	UserID      int64

	Ctime int64
	Utime int64
}

func NewAccountDao(db *gorm.DB) *AccountDAO {
	return &AccountDAO{
		db: db,
	}
}

func (a *AccountDAO) GetAccount(ctx *gin.Context, userId int64) ([]Account, error) {
	var accounts []Account
	err := a.db.WithContext(ctx).Where("user_id = ?", userId).Find(&accounts).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return accounts, err
}

func (a *AccountDAO) GetAccountWithType(ctx *gin.Context, userId int64, accountType string) (Account, error) {
	var account Account
	err := a.db.WithContext(ctx).Where("user_id = ? and account_type = ?", userId, accountType).Find(&account).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return account, err
}

func (a *AccountDAO) CreateOrUpdate(ctx *gin.Context, data Account) (Account, error) {
	where := Account{
		AccountType: cast.ToString(data.AccountType),
		UserID:      cast.ToInt64(data.UserID),
	}
	var account Account
	err := a.db.Where(where).Assign(data).FirstOrCreate(&account).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return account, err
}
