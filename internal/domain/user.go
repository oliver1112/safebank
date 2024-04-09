package domain

// User domain, is entity of DDD, or BO(business object)
type User struct {
	ID       int64
	Email    string
	Password string
}
