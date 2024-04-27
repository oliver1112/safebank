package domain

type Response struct {
	Status   int64
	Data     interface{}
	ErrorMsg string
}

type AccountData struct {
	SavingAccount       map[string]interface{} `json:"saving_account"`
	CheckingAccount     map[string]interface{} `json:"checking_account"`
	PersonalLoanAccount map[string]interface{} `json:"personal_loan_account"`
	HomeLoanAccount     map[string]interface{} `json:"home_loan_account"`
	StudentLoanAccount  map[string]interface{} `json:"student_loan_account"`
}

type UserCenter struct {
	UserInfo    User        `json:"user_info"`
	AccountInfo AccountData `json:"account_info"`
}
