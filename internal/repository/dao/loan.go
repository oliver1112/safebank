package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanDAO struct {
	db *gorm.DB
}

type Loan struct {
	AccountID int64 `gorm:"primaryKey" json:"account_id"`
	Account   Account
	Rate      float64 `json:"rate"`
	Amount    float64 `json:"amount"`
	Month     int     `json:"month"`
	Payment   float64 `json:"payment"`
	Type      string  `json:"type"`
}

func NewLoanDao(db *gorm.DB) *LoanDAO {
	return &LoanDAO{
		db: db,
	}
}

func (sd *LoanDAO) GetLoan(ctx *gin.Context, accountId int64) (Loan, error) {
	var loan Loan
	err := sd.db.WithContext(ctx).Where("account_id = ?", accountId).First(&loan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return loan, err
}

func (sd *LoanDAO) GetByAccountID(ctx *gin.Context, accountID int64) (Loan, error) {
	var loan Loan
	err := sd.db.WithContext(ctx).Where("account_id = ?", accountID).First(&loan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return loan, err
}

func (sd *LoanDAO) CreateOrUpdate(ctx *gin.Context, data Loan) (Loan, error) {
	where := Loan{
		AccountID: data.AccountID,
	}
	var loan Loan
	err := sd.db.Where(where).Assign(data).FirstOrCreate(&loan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return loan, err
}
