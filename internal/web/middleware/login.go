package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"safebank/internal/lib"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// no need to login and check
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}

		data, err := ctx.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("data: %v\n", string(data))

		m := map[string]string{}
		_ = json.Unmarshal(data, &m)
		fmt.Printf("id: %d\n", m["id"])

		// rewrite data to body
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		userToken := lib.UserToken{}
		userToken.DecodeToken(m["userToken"])

		id := userToken.UserID
		if id <= 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("userID", id)
	}
}
