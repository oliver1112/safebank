package domain

type Account struct {
	ID   int64
	name string

	Street      string
	City        string
	State       string
	Zip         string
	AccountType string
	Apart       string
	UserID      int64

	Ctime int64
	Utime int64
}
