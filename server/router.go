package server

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"singo/api"
	"singo/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", api.Ping)

		user := v1.Group("user")

		// 用户登录
		user.POST("login", api.UserLogin)

		// 用户注册
		user.POST("register", api.UserRegister)

		user.Use(middleware.AuthMiddleware())

		// 需要登录保护的
		user.GET("info", api.UserMe)

		user.GET("list", api.Get)
	}

	{
		player := v1.Group("player")

		// 用户登录
		player.GET(":code", api.PlayerLogin)

		player.Use(middleware.AuthMiddleware())

		player.POST("start", api.StartGame)

		player.GET("energy", api.MyEnergy)

		player.PUT("health", api.AddHealth)

		player.POST("item", api.AddItem)

		player.POST("share", api.Share)
		player.GET("team", api.MyTeam)
		player.GET("result", api.MyResult)
		player.GET("coin", api.MyCoin)
		player.POST("coin", api.PostCoin)
		player.POST("team", api.PostTeam)
		player.POST("result", api.PostResult)
	}

	return r
}
