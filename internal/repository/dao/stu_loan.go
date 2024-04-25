package dao

import "gorm.io/gorm"

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
