package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HomeLoanDAO struct {
	db *gorm.DB
}

type HomeLoan struct {
	LoanID int64 `gorm:"primaryKey"`
	Loan   Loan  //`gorm:"foreignKey:AccountID"`

	BuildYear    int
	InsurAccNum  int
	InsurName    string
	InsurStreet  string
	InsurCity    string
	InsurState   string
	InsurZip     int
	YearInsurPrm float64
}

func NewHomeLoanDao(db *gorm.DB) *HomeLoanDAO {
	return &HomeLoanDAO{
		db: db,
	}
}

func (hd *HomeLoanDAO) GetHome(ctx *gin.Context, userId int64) (HomeLoan, error) {
	var homeLoan HomeLoan
	err := hd.db.WithContext(ctx).Where("user_id = ?", userId).First(&homeLoan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return homeLoan, err
}

func (hd *HomeLoanDAO) CreateOrUpdate(ctx *gin.Context, data HomeLoan) (HomeLoan, error) {
	where := HomeLoan{
		LoanID: data.LoanID,
	}
	var homeLoan HomeLoan
	err := hd.db.Where(where).Assign(data).FirstOrCreate(&homeLoan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return homeLoan, err
}
