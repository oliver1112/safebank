package domain

type Saving struct {
	AccountID    int
	Account      Account
	InterestRate float64
	Amount       int64
}
