package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"safebank/internal/lib"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// no need to login and check
		if strings.HasPrefix(ctx.Request.URL.Path, "/admin") ||
			ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}

		data, err := ctx.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("data: %v\n", string(data))

		m := make(map[string]interface{})
		_ = json.Unmarshal(data, &m)
		fmt.Printf("userToken: %s\n", cast.ToString(m["userToken"]))

		// rewrite data to body
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		userToken := lib.UserToken{}
		userToken.DecodeToken(cast.ToString(m["userToken"]))
		id := userToken.UserID

		if id <= 0 || userToken.ExpiresAt < time.Now().Unix() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Printf("userID: %d\n", id)
		ctx.Set("userID", id)
	}
}
