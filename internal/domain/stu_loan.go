package domain

type StuLoan struct {
	LoanID int64
	Loan   Loan

	EduInstitute    string
	StudentID       int
	GradStatus      string
	ExpectGradMonth int
	ExpectGradYear  int
}
