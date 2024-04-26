package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"safebank/internal/lib"
	"strings"
)

type AdminLoginMiddlewareBuilder struct {
}

func NewAdminLoginMiddlewareBuilder() *AdminLoginMiddlewareBuilder {
	return &AdminLoginMiddlewareBuilder{}
}

func (a *AdminLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// no need to login and check
		if !strings.HasPrefix(ctx.Request.URL.Path, "/admin") ||
			ctx.Request.URL.Path == "/admin/login" {
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
		userToken.DecodeToken(m["adminToken"])

		id := userToken.UserID
		if id <= 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("adminID", id)
	}
}
