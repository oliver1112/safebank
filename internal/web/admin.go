package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"safebank/internal/domain"
	"safebank/internal/lib"
	"safebank/internal/repository/dao"
	"safebank/internal/service"
	"time"
)

type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	svc := service.NewAdminService(db)
	// init an employee account(Assumption it was added by a manager)
	hash, _ := bcrypt.GenerateFromPassword([]byte("aaaAAA@@@111"), bcrypt.DefaultCost)
	employeeAssumption := dao.Employee{
		Email:    "wsj@gmail.com",
		Password: string(hash),
		FName:    "WSJ",
		LName:    "Group",
	}

	employee, _ := svc.EmployeeDao.CreateOrUpdate(employeeAssumption)
	fmt.Println("init admin user")
	fmt.Printf("%v", employee)
	return &AdminHandler{
		svc: svc,
	}
}

func (a *AdminHandler) RegisterRoutes(server *gin.Engine) {
	adminRouterG := server.Group("/admin")
	adminRouterG.POST("/login", a.Login)
	adminRouterG.POST("/getaccountsbyemail", a.GetAccountsByEmail)
	adminRouterG.POST("/getaccountbyaccountid", a.GetAccountsByAccountID)
	adminRouterG.POST("/updateaccount", a.UpdateAccount)
	adminRouterG.POST("/dashboard", a.DashBoard)
	adminRouterG.POST("/deleteaccount", a.DeleteAccount)
}

func (a *AdminHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type ResponseData struct {
		AdminToken string `json:"adminToken"`
	}
	responseData := ResponseData{}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "args error",
			Data:     responseData,
		})
		return
	}

	user, err := a.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "wrong email or password",
			Data:     responseData,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   2,
			ErrorMsg: "system error",
			Data:     responseData,
		})
		return
	}

	// Login success and set session
	userToken := lib.UserToken{
		UserID:    user.ID,
		ExpiresAt: time.Now().Unix() + 3600*24,
		IP:        ctx.ClientIP(),
	}

	responseData.AdminToken = userToken.EncodeToken()
	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     responseData,
	})

	return
}

func (a *AdminHandler) GetAccountsByAccountID(ctx *gin.Context) {
	type Req struct {
		AccountID int64 `json:"account_id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "args error",
		})
		return
	}

	accountData, err := a.svc.AccountService.GetAccountsByAccountID(ctx, req.AccountID)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: cast.ToString(err),
			Data:     accountData,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     accountData,
	})
	return
}

func (a *AdminHandler) GetAccountsByEmail(ctx *gin.Context) {
	type Req struct {
		Email string `json:"email"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "args error",
		})
		return
	}

	accountData, err := a.svc.AccountService.GetAccountsByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: cast.ToString(err),
			Data:     accountData,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     accountData,
	})
	return
}

func (a *AdminHandler) UpdateAccount(ctx *gin.Context) {
	type Req struct {
		AccountID  int64                  `json:"account_id"`
		UpdateData map[string]interface{} `json:"update_data"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "args error",
		})
		return
	}

	isUpdate, err := a.svc.UpdateAccountInfo(ctx, req.AccountID, req.UpdateData)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: cast.ToString(err),
			Data:     false,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     isUpdate,
	})
	return
}

func (a *AdminHandler) DashBoard(ctx *gin.Context) {
	response := make(map[string]interface{})
	var totalUserNum int64
	var savingAccountNum int64
	var checkingAccountNum int64
	var homeLoanAccountNum int64
	var studentLoanAccountNum int64
	var personalLoanAccountNum int64
	type TotalSavingDeposit struct {
		Total int64
	}
	type TotalCheckingDeposit struct {
		Total int64
	}
	type TotalLoanAmount struct {
		Total float64
	}
	var totalSavingDeposit TotalSavingDeposit
	var totalCheckingDeposit TotalCheckingDeposit
	var totalLoanAmount TotalLoanAmount

	a.svc.AccountDao.Db.Model(&dao.User{}).Count(&totalUserNum)
	a.svc.AccountDao.Db.Model(&dao.Account{}).Where(&dao.Account{AccountType: "S"}).Count(&savingAccountNum)
	a.svc.AccountDao.Db.Model(&dao.Account{}).Where(&dao.Account{AccountType: "C"}).Count(&checkingAccountNum)
	a.svc.LoanDao.Db.Model(&dao.Loan{}).Where(&dao.Loan{Type: "H"}).Count(&homeLoanAccountNum)
	a.svc.LoanDao.Db.Model(&dao.Loan{}).Where(&dao.Loan{Type: "S"}).Count(&studentLoanAccountNum)
	a.svc.LoanDao.Db.Model(&dao.Loan{}).Where(&dao.Loan{Type: "L"}).Count(&personalLoanAccountNum)
	a.svc.SavingDao.Db.Model(&dao.Saving{}).Select("sum(amount) as total").Find(&totalSavingDeposit)
	a.svc.CheckingDao.Db.Model(&dao.Checking{}).Select("sum(amount) as total").Find(&totalCheckingDeposit)
	a.svc.LoanDao.Db.Model(&dao.Loan{}).Select("sum(amount) as total").Find(&totalLoanAmount)

	response["saving_account_num"] = savingAccountNum
	response["checking_account_num"] = checkingAccountNum
	response["home_loan_account_num"] = homeLoanAccountNum
	response["student_loan_account_num"] = studentLoanAccountNum
	response["personal_loan_account_num"] = personalLoanAccountNum
	response["total_saving_deposit"] = cast.ToFloat64(totalSavingDeposit.Total) / 100
	response["total_checking_deposit"] = cast.ToFloat64(totalCheckingDeposit.Total) / 100
	response["total_loan_amount"] = totalLoanAmount.Total

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     response,
	})
	return
}

func (a *AdminHandler) DeleteAccount(ctx *gin.Context) {
	type Req struct {
		AccountID int64 `json:"account_id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "args error",
		})
		return
	}

	a.svc.DeleteAccountByID(ctx, req.AccountID)

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
	})
	return
}
