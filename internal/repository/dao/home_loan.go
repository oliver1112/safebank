package dao

import "gorm.io/gorm"

type HomeLoanDAO struct {
	db *gorm.DB
}

type HomeLoan struct {
	AccountID int64 `gorm:"primaryKey"`
	Loan      Loan  `gorm:"foreignKey:AccountID"`

	BuildYear    int
	InsurAccNum  int
	InsurName    string
	InsurStreet  string
	InsurCity    string
	InsurState   string
	InsurZip     int
	YearInsurPrm float64
}
