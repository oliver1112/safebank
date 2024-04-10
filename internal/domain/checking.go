package domain

type Checking struct {
	AccountID     int64
	Account       Account
	ServiceCharge float64
	Amount        float64
}
