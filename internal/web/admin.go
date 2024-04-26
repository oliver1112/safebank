package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"safebank/internal/domain"
	"safebank/internal/lib"
	"safebank/internal/repository/dao"
	"safebank/internal/service"
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
		UserID: user.ID,
	}

	responseData.AdminToken = userToken.EncodeToken()
	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     responseData,
	})

	return
}
