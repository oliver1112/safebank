package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionDAO struct {
	Db *gorm.DB
}

func NewTransactionDao(db *gorm.DB) *TransactionDAO {
	return &TransactionDAO{
		Db: db,
	}
}

type Transaction struct {
	TransactionID   int64  `gorm:"primaryKey" json:"Transaction_id"`
	FromAccountID   int64  `gorm:"index" json:"from_account_id"`
	FromAccountName string `gorm:"index" json:"from_account_name"`
	ToAccountID     int64  `gorm:"index" json:"to_account_id"`
	ToAccountName   string `json:"to_account_name"`
	Amount          int64  `json:"amount"`
}

func (Transaction) TableName() string {

	return "wsj_transactions"
}

func (i *TransactionDAO) GetByAccountID(ctx *gin.Context, AccountID int64) ([]Transaction, error) {
	var transactions []Transaction
	err := i.Db.WithContext(ctx).Where(&Transaction{FromAccountID: AccountID}).Or(&Transaction{ToAccountID: AccountID}).Find(&transactions).Error
	return transactions, err
}

func (i *TransactionDAO) Create(ctx *gin.Context, data Transaction) error {
	err := i.Db.Create(&data).Error
	return err
}
