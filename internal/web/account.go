package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"safebank/internal/domain"
	"safebank/internal/repository/dao"
	"safebank/internal/service"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {

	svc := service.NewAccountService(db)

	return &AccountHandler{
		svc: svc,
	}
}

func (a *AccountHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/account")
	ug.POST("/test", a.CreateOrUpdateSavingAccount)

	ug.POST("/addsaving", a.CreateOrUpdateSavingAccount)
	ug.POST("/addchecking", a.CreateOrUpdateCheckingAccount)
	ug.POST("/addloan", a.CreateOrUpdateLoan)
	ug.POST("/addhomeloan", a.CreateOrUpdateHomeLoan)
	ug.POST("/addstudentloan", a.CreateOrUpdateStuLoan)

	ug.POST("/usercenter", a.UserCenter)
	//ug.GET("/checking", a.FindCheckingAccount)
	//ug.GET("/loan", a.FindLoan)
	//ug.GET("/homeloan", a.FindHomeLoan)
	//ug.GET("/studentloan", a.FindStuLoan)
}

func (a *AccountHandler) UserCenter(ctx *gin.Context) {

	userId, _ := ctx.Get("userID")
	userID := cast.ToInt64(userId)

	userInfo, _ := a.svc.UserRepo.FindUserByUserID(ctx, userID)

	fmt.Printf("%v", userInfo)

	accountList, _ := a.svc.AccountDao.GetAccountList(ctx, userID)

	type accountData struct {
		Saving_account        dao.Saving
		Checking_account      dao.Checking
		Personal_loan_account dao.Loan
		Home_loan_account     dao.HomeLoan
		Student_loan_account  dao.StuLoan
	}

	var data accountData

	for _, account := range accountList {
		if account.AccountType == "S" {
			data.Saving_account, _ = a.svc.SavingDao.GetSaving(ctx, account.ID)
			data.Saving_account.Account = account

		} else if account.AccountType == "C" {
			data.Checking_account, _ = a.svc.CheckingDao.GetChecking(ctx, account.ID)
			data.Checking_account.Account = account

		} else if account.AccountType == "L" {
			loanData, _ := a.svc.LoanDao.GetLoan(ctx, account.ID)

			if loanData.Type == "L" {
				data.Personal_loan_account = loanData
				data.Personal_loan_account.Account = account

			} else if loanData.Type == "H" {
				data.Home_loan_account, _ = a.svc.HomeLoanDao.GetHomeLoan(ctx, account.ID)
				data.Home_loan_account.Loan = loanData
				data.Home_loan_account.Loan.Account = account

			} else if loanData.Type == "S" {
				data.Student_loan_account, _ = a.svc.StuLoanDao.GetStuLoan(ctx, account.ID)
				data.Student_loan_account.Loan = loanData
				data.Student_loan_account.Loan.Account = account
			}
		}
	}

	type Info struct {
		UserInfo    domain.User
		AccountInfo accountData
	}

	var info Info
	info.UserInfo = userInfo
	info.AccountInfo = data

	infoJSON, _ := json.MarshalIndent(info, "", "    ")

	fmt.Printf("\n%v\n", string(infoJSON))
	//ctx.JSON(http.StatusOK, string(infoJSON))
	ctx.Data(http.StatusOK, "application/json", infoJSON)
}

func (a *AccountHandler) CreateOrUpdateSavingAccount(ctx *gin.Context) {
	type Req struct {
		Street  string `json:"address"`
		Apart   string `json:"address2"`
		Country string `json:"country"`
		State   string `json:"state"`
		City    string `json:"city"`
		Zip     string `json:"zip"`
	}
	var responseData interface{}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "Args error",
			Data:     responseData,
		})
		return
	}

	userId, _ := ctx.Get("userID")

	data := dao.Account{
		Name:        "SavingAccount" + cast.ToString(rand.Intn(9999999)+1000000),
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Zip:         req.Zip,
		Apart:       req.Apart,
		AccountType: "S",
		UserID:      cast.ToInt64(userId),
	}

	account, err := a.svc.AccountDao.CreateOrUpdate(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	randomValue := 100.00 + rand.Float64()*100.00
	// Truncate to two decimal places
	randomRate := float64(int(randomValue*100)) / 100

	savingData := dao.Saving{
		AccountID:    account.ID,
		InterestRate: randomRate,
		Amount:       0,
		Account:      account,
	}

	saving, err := a.svc.SavingDao.CreateOrUpdate(ctx, savingData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     saving,
	})
}

func (a *AccountHandler) CreateOrUpdateCheckingAccount(ctx *gin.Context) {
	type Req struct {
		Street  string `json:"address"`
		Apart   string `json:"address2"`
		Country string `json:"country"`
		State   string `json:"state"`
		City    string `json:"city"`
		Zip     string `json:"zip"`
	}

	var req Req
	var responseData interface{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "Args error",
			Data:     responseData,
		})
		return
	}

	userId, _ := ctx.Get("userID")

	data := dao.Account{
		Name:        "CheckingAccount" + cast.ToString(rand.Intn(9999999)+1000000),
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Zip:         req.Zip,
		Apart:       req.Apart,
		AccountType: "C",
		UserID:      cast.ToInt64(userId),
	}

	account, err := a.svc.AccountDao.CreateOrUpdate(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	// checking account
	randomValue := rand.Float64() * 10.00
	// Truncate to two decimal places
	randomRate := float64(int(randomValue*100)) / 100

	checkingData := dao.Checking{
		AccountID:     account.ID,
		ServiceCharge: randomRate,
		Amount:        0,
		Account:       account,
	}

	saving, err := a.svc.CheckingDao.CreateOrUpdate(ctx, checkingData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     saving,
	})
}

func (a *AccountHandler) CreateOrUpdateLoan(ctx *gin.Context) {
	type Req struct {
		loanAmount float64 `json:"lamount"`
		loanMonth  int     `json:"lmonth"`
		Street     string  `json:"address"`
		Apart      string  `json:"address2"`
		Country    string  `json:"country"`
		State      string  `json:"state"`
		City       string  `json:"city"`
		Zip        string  `json:"zip"`
	}

	var req Req
	var responseData interface{}

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "Args error",
			Data:     responseData,
		})
		return
	}

	userId, _ := ctx.Get("userID")

	data := dao.Account{
		Name:        "PersonalLoanAccount" + cast.ToString(rand.Intn(9999999)+1000000),
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Zip:         req.Zip,
		Apart:       req.Apart,
		AccountType: "L",
		UserID:      cast.ToInt64(userId),
	}

	account, err := a.svc.AccountDao.CreateOrUpdate(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	randomValue := rand.Float64() * 50.00
	// Truncate to two decimal places
	randomRate := float64(int(randomValue*100)) / 100

	loanData := dao.Loan{
		AccountID: account.ID,
		Rate:      randomRate,
		Amount:    req.loanAmount,
		Account:   account,
		Month:     req.loanMonth,
		Payment:   0,
		Type:      "L",
	}

	loan, err := a.svc.LoanDao.CreateOrUpdate(ctx, loanData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     loan,
	})
}

func (a *AccountHandler) CreateOrUpdateHomeLoan(ctx *gin.Context) {
	type Req struct {
		loanAmount   float64 `json:"lamount"`
		loanMonth    int     `json:"lmonth"`
		buildYear    string  `json:"buildYear"`
		InsurAccNum  string  `json:"insu_acc_no"`
		InsurName    string  `json:"insu_name"`
		InsurCountry string  `json:"insu_country"`
		InsurStreet  string  `json:"insu_address"`
		InsurCity    string  `json:"insu_city"`
		InsurState   string  `json:"insu_state"`
		InsurZip     string  `json:"insu_zip"`
		Street       string  `json:"address"`
		Apart        string  `json:"address2"`
		Country      string  `json:"country"`
		State        string  `json:"state"`
		City         string  `json:"city"`
		Zip          string  `json:"zip"`
	}

	var req Req
	var responseData interface{}

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "Args error",
			Data:     responseData,
		})
		return
	}

	userId, _ := ctx.Get("userID")

	data := dao.Account{
		Name:        "HomeLoanAccount" + cast.ToString(rand.Intn(9999999)+1000000),
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Zip:         req.Zip,
		Apart:       req.Apart,
		AccountType: "L",
		UserID:      cast.ToInt64(userId),
	}

	account, err := a.svc.AccountDao.CreateOrUpdate(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	randomValue := rand.Float64() * 50.00
	// Truncate to two decimal places
	randomRate := float64(int(randomValue*100)) / 100

	loanData := dao.Loan{
		AccountID: account.ID,
		Rate:      randomRate,
		Amount:    req.loanAmount,
		Account:   account,
		Month:     req.loanMonth,
		Payment:   0,
		Type:      "H",
	}

	loan, err := a.svc.LoanDao.CreateOrUpdate(ctx, loanData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	randomValue = rand.Float64() * 20.00
	// Truncate to two decimal places
	randomPrm := float64(int(randomValue*100)) / 100

	homeLoanData := dao.HomeLoan{
		Loan:         loan,
		BuildYear:    cast.ToInt(req.buildYear),
		InsurAccNum:  cast.ToInt(req.InsurAccNum),
		InsurName:    req.InsurName,
		InsurStreet:  req.InsurStreet,
		InsurCity:    req.InsurCity,
		InsurState:   req.InsurState,
		InsurZip:     cast.ToInt(req.InsurZip),
		YearInsurPrm: randomPrm,
	}

	homeLoan, err := a.svc.HomeLoanDao.CreateOrUpdate(ctx, homeLoanData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     homeLoan,
	})
}

func (a *AccountHandler) CreateOrUpdateStuLoan(ctx *gin.Context) {
	type Req struct {
		loanAmount      float64 `json:"lamount"`
		loanMonth       int     `json:"lmonth"`
		EduInstitute    string  `json:"eduinsititute"`
		StudentID       string  `json:"studentid"`
		GradStatus      string  `json:"grad_status"`
		ExpectGradMonth string  `json:"graduationMonth"`
		ExpectGradYear  string  `json:"graduationYear"`
		Street          string  `json:"address"`
		Apart           string  `json:"address2"`
		Country         string  `json:"country"`
		State           string  `json:"state"`
		City            string  `json:"city"`
		Zip             string  `json:"zip"`
	}

	var req Req
	var responseData interface{}

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "Args error",
			Data:     responseData,
		})
		return
	}

	userId, _ := ctx.Get("userID")

	data := dao.Account{
		Name:        "StudentLoanAccount" + cast.ToString(rand.Intn(9999999)+1000000),
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Zip:         req.Zip,
		Apart:       req.Apart,
		AccountType: "L",
		UserID:      cast.ToInt64(userId),
	}

	account, err := a.svc.AccountDao.CreateOrUpdate(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	randomValue := rand.Float64() * 50.00
	// Truncate to two decimal places
	randomRate := float64(int(randomValue*100)) / 100

	loanData := dao.Loan{
		AccountID: account.ID,
		Rate:      randomRate,
		Amount:    req.loanAmount,
		Account:   account,
		Month:     req.loanMonth,
		Payment:   0,
		Type:      "S",
	}

	loan, err := a.svc.LoanDao.CreateOrUpdate(ctx, loanData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	stuLoanData := dao.StuLoan{
		Loan:            loan,
		EduInstitute:    req.EduInstitute,
		StudentID:       cast.ToInt(req.StudentID),
		GradStatus:      req.GradStatus,
		ExpectGradMonth: cast.ToInt(req.ExpectGradMonth),
		ExpectGradYear:  cast.ToInt(req.ExpectGradYear),
	}

	stuLoan, err := a.svc.StuLoanDao.CreateOrUpdate(ctx, stuLoanData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "db error",
			Data:     responseData,
		})
	}

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     stuLoan,
	})
}
