package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StuLoanDAO struct {
	db *gorm.DB
}

func NewStuLoanDao(db *gorm.DB) *StuLoanDAO {
	return &StuLoanDAO{
		db: db,
	}
}

type StuLoan struct {
	AccountID int64 `json:"account_id"`
	Loan      Loan  `gorm:"foreignKey:AccountID"`

	EduInstitute    string `json:"edu_institute"`
	StudentID       int    `json:"student_id"`
	GradStatus      string `json:"grad_status"`
	ExpectGradMonth int    `json:"expect_grad_month"`
	ExpectGradYear  int    `json:"expect_grad_year"`
}

func (sd *StuLoanDAO) GetStuLoan(ctx *gin.Context, accountID int64) (StuLoan, error) {
	var stuLoan StuLoan
	err := sd.db.WithContext(ctx).Where("account_id = ?", accountID).First(&stuLoan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return stuLoan, err
}

func (sd *StuLoanDAO) CreateOrUpdate(ctx *gin.Context, data StuLoan) (StuLoan, error) {
	where := StuLoan{
		AccountID: data.AccountID,
	}
	var stuLoan StuLoan
	err := sd.db.Where(where).Assign(data).FirstOrCreate(&stuLoan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return stuLoan, err
}
