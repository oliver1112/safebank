package domain

// User domain, is entity of DDD, or BO(business object)
type User struct {
	ID    int64
	Email string

	FName   string
	LName   string
	Country string

	State  string
	City   string
	Street string
	Apart  string
	Zip    string

	Password string
	Ctime    int64
	Utime    int64

	Account []Account
}
