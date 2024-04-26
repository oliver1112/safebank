package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"safebank/internal/repository"
	"safebank/internal/repository/dao"
	"safebank/internal/service"
	"safebank/internal/web"
	"safebank/internal/web/middleware"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	accountHandler := initAccount(db)
	accountHandler.RegisterRoutes(server)

	adminHandler := initAdmin(db)
	adminHandler.RegisterRoutes(server)

	server.Run("0.0.0.0:8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// CORS (Cross-Origin Resource Sharing)
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			//return strings.Contains(origin, "yourcompany.com")
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().Build())
	server.Use(middleware.NewAdminLoginMiddlewareBuilder().Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initAccount(db *gorm.DB) *web.AccountHandler {
	return web.NewAccountHandler(db)
}

func initAdmin(db *gorm.DB) *web.AdminHandler {
	return web.NewAdminHandler(db)
}

func initDB() *gorm.DB {
	//change localhost to 43.130.62.214
	db, err := gorm.Open(mysql.Open("root:root@tcp(43.130.62.214:13333)/safebank"))
	//db, err := gorm.Open(mysql.Open("root:Morjoe123@tcp(localhost:3306)/safebank"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
