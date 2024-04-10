package domain

type Saving struct {
	AccountID    int64
	Account      Account
	InterestRate float64
	Amount       int64
}
