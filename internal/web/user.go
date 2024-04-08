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

	//receive request
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//check user signup valid logic
	isEmail, err := u.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "invalid email")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "the password is not the same")
		return
	}

	isPassword, err := u.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "Password must contain letters, numbers, special characters, and be no less than eight characters")
		return
	}

	// call service signup
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "email conflict")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "system error")
	}

	ctx.String(http.StatusOK, "register success")
	fmt.Printf("%v", req)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "wrong email or password")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	// Login success and set session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "Login Success")

	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

	ctx.String(http.StatusOK, "This is your profile")
}
