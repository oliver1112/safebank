package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"safebank/internal/domain"
	"safebank/internal/service"
)

const (
	// regular expression of email and password
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	svc            *service.UserService
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:            svc,
		emailRexExp:    emailExp,
		passwordRexExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", u.SignUp)
	// POST /users/login
	ug.POST("/login", u.Login)
	// POST /users/edit
	ug.POST("/edit", u.Edit)
	// GET /users/profile
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var responseData interface{}

	//receive request
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   -1,
			ErrorMsg: "args error",
			Data:     responseData,
		})
		return
	}

	//check user signup valid logic
	isEmail, err := u.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   1,
			ErrorMsg: "system error",
			Data:     responseData,
		})
		return
	}
	if !isEmail {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   2,
			ErrorMsg: "invalid email",
			Data:     responseData,
		})
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   3,
			ErrorMsg: "the password is not the same",
			Data:     responseData,
		})
		return
	}

	isPassword, err := u.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   4,
			ErrorMsg: "system error",
			Data:     responseData,
		})
		return
	}
	if !isPassword {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   5,
			ErrorMsg: "Password must contain letters, numbers, special characters, and be no less than eight characters",
			Data:     responseData,
		})
		return
	}

	// call service signup
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   6,
			ErrorMsg: "email conflict",
			Data:     responseData,
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, domain.Response{
			Status:   7,
			ErrorMsg: "system error",
			Data:     responseData,
		})
	}

	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     responseData,
	})
	fmt.Printf("%v", req)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var responseData interface{}
	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
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
	sess := sessions.Default(ctx)
	sess.Set("userId", user.ID)
	sess.Save()
	ctx.JSON(http.StatusOK, domain.Response{
		Status:   0,
		ErrorMsg: "",
		Data:     responseData,
	})

	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

	sess := sessions.Default(ctx)
	id := sess.Get("userId")

	if id == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, ok := id.(int64)
	if !ok {
		ctx.String(http.StatusOK, "session error")
		return
	}

	//println(useID)
	userData, err := u.svc.GetProfile(userID)
	if err != nil {
		ctx.String(http.StatusOK, "no such user")
		return
	}

	responseData := gin.H{
		"status": 0,
		// "data": userData
		"data":   userData,
		"errmsg": "",
	}

	// Return JSON response
	ctx.JSON(200, responseData)
}

func (u *UserHandler) AddSaving(ctx *gin.Context) {
	type SavingReq struct {
		Email string `json:""`
	}

	var req SavingReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	sess := sessions.Default(ctx)
	id := sess.Get("userId")

	if id == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, ok := id.(int64)
	print(userID)
	if !ok {
		ctx.String(http.StatusOK, "session error")
		return
	}

}
