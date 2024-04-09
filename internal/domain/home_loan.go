package domain

type HomeLoan struct {
	LoanID int64 //`gorm:"primaryKey"`
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
