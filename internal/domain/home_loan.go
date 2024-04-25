package domain

type HomeLoan struct {
	LoanID int64
	Loan   Loan

	BuildYear    int
	InsurAccNum  int
	InsurName    string
	InsurCountry string
	InsurStreet  string
	InsurCity    string
	InsurState   string
	InsurZip     int
	YearInsurPrm float64
}
