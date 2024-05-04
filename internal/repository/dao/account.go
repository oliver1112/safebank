package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type AccountDAO struct {
	Db *gorm.DB
}

type Account struct {
	ID          int64  `gorm:"primaryKey,autoIncrement" json:"id"`
	Name        string `json:"name"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Apart       string `json:"apart"`
	AccountType string `json:"account_type"`
	UserID      int64  `json:"user_id"`

	Ctime int64 `gorm:"autoUpdateTime" json:"ctime"`
	Utime int64 `gorm:"autoCreateTime" json:"utime"`
}

func (Account) TableName() string {
	return "wsj_account"
}

func NewAccountDao(db *gorm.DB) *AccountDAO {
	return &AccountDAO{
		Db: db,
	}
}

func (a *AccountDAO) GetAccountList(ctx *gin.Context, userId int64) ([]Account, error) {
	var accounts []Account
	err := a.Db.WithContext(ctx).Where("user_id = ?", userId).Find(&accounts).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return accounts, err
}

func (a *AccountDAO) GetAccountByID(ctx *gin.Context, accountID int64) (Account, error) {
	var account Account
	err := a.Db.WithContext(ctx).Where("id = ?", accountID).First(&account).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return account, err
}

func (a *AccountDAO) GetAccountWithType(ctx *gin.Context, userId int64, accountType string) (Account, error) {
	var account Account
	err := a.Db.WithContext(ctx).Where("user_id = ? and account_type = ?", userId, accountType).Find(&account).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return account, err
}

func (a *AccountDAO) CreateOrUpdate(ctx *gin.Context, data Account) (Account, error) {
	where := Account{
		AccountType: cast.ToString(data.AccountType),
		UserID:      cast.ToInt64(data.UserID),
	}
	var account Account
	err := a.Db.Where(where).Assign(data).FirstOrCreate(&account).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return account, err
}
