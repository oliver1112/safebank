package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HomeLoanDAO struct {
	Db *gorm.DB
}

type HomeLoan struct {
	AccountID int64 `json:"account_id"`
	Loan      Loan  `gorm:"foreignKey:AccountID"`

	BuildYear    int     `json:"build_year"`
	InsurAccNum  int     `json:"insur_acc_num"`
	InsurName    string  `json:"insur_name"`
	InsurStreet  string  `json:"insur_street"`
	InsurCity    string  `json:"insur_city"`
	InsurState   string  `json:"insur_state"`
	InsurZip     int     `json:"insur_zip"`
	YearInsurPrm float64 `json:"year_insur_prm"`
}

func (HomeLoan) TableName() string {
	return "wsj_homeloan"
}

func NewHomeLoanDao(db *gorm.DB) *HomeLoanDAO {
	return &HomeLoanDAO{
		Db: db,
	}
}

func (hd *HomeLoanDAO) GetHomeLoan(ctx *gin.Context, accountId int64) (HomeLoan, error) {
	var homeLoan HomeLoan
	err := hd.Db.WithContext(ctx).Where("account_id = ?", accountId).First(&homeLoan).Error
	return homeLoan, err
}

func (hd *HomeLoanDAO) CreateOrUpdate(ctx *gin.Context, data HomeLoan) (HomeLoan, error) {
	where := HomeLoan{
		AccountID: data.AccountID,
	}
	var homeLoan HomeLoan
	err := hd.Db.Where(where).Assign(data).FirstOrCreate(&homeLoan).Error
	return homeLoan, err
}
