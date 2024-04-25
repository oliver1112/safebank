package web

import (
	"fmt"
	"github.com/gin-contrib/sessions"
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

	session := sessions.Default(ctx)
	id := session.Get("userId")
	//accounts, err := a.svc.GetAccount(ctx)
	//if err != nil {
	//	ctx.String(http.StatusOK, "system error")
	//	return
	//}
	//
	//if len(accounts) == 0 {
	//	// create
	//
	//}

	data := dao.Account{
		Name:        "SavingAccount" + cast.ToString(rand.Intn(9999999)+1000000),
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Zip:         req.Zip,
		Apart:       req.Apart,
		AccountType: "S",
		UserID:      cast.ToInt64(id),
	}

	account, err := a.svc.AccountDao.CreateOrUpdate(ctx, data)
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
		Data:     account,
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
	if err := ctx.Bind(&req); err != nil {
		return
	}

	fmt.Printf("%v", req)
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
	if err := ctx.Bind(&req); err != nil {
		return
	}

	fmt.Printf("%v", req)
}

func (a *AccountHandler) CreateOrUpdateStuLoan(ctx *gin.Context) {
	type Req struct {
		loanAmount      float64 `json:"lamount"`
		loanMonth       int     `json:"lmonth"`
		EduInstitute    string  `json:"eduinsititute"`
		StudentID       string  `json:"sid"`
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
	if err := ctx.Bind(&req); err != nil {
		return
	}

	fmt.Printf("%v", req)
}
