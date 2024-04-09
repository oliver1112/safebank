package dao

import "gorm.io/gorm"

type StuLoanDAO struct {
	db *gorm.DB
}

type StuLoan struct {
	LoanID int64 `gorm:"primaryKey"`
	Loan   Loan

	EduInstitute    string
	StudentID       int
	GradStatus      string
	ExpectGradMonth int
	ExpectGradYear  int
}