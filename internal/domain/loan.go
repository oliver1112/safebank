package domain

type Loan struct {
	AccountID int64
	Account   Account
	Rate      float64
	Amount    float64
	Month     int
	Payment   float64
	Type      string
}
