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
	LoanID int64 `gorm:"primaryKey"`
	Loan   Loan  //`gorm:"foreignKey:AccountID"`

	EduInstitute    string
	StudentID       int
	GradStatus      string
	ExpectGradMonth int
	ExpectGradYear  int
}

func (sd *StuLoanDAO) GetStuLoan(ctx *gin.Context, userId int64) (StuLoan, error) {
	var stuLoan StuLoan
	err := sd.db.WithContext(ctx).Where("loan_id = ?", userId).First(&stuLoan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return stuLoan, err
}

func (sd *StuLoanDAO) CreateOrUpdate(ctx *gin.Context, data StuLoan) (StuLoan, error) {
	where := StuLoan{
		LoanID: data.LoanID,
	}
	var stuLoan StuLoan
	err := sd.db.Where(where).Assign(data).FirstOrCreate(&stuLoan).Error
	//err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return stuLoan, err
}
