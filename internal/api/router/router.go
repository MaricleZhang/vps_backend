package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mariclezhang/vps_backend/internal/api/handler"
	"github.com/mariclezhang/vps_backend/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter(frontendURL string) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware(frontendURL))
	r.Use(middleware.LoggerMiddleware())

	// 初始化处理器
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	subscriptionHandler := handler.NewSubscriptionHandler()
	nodeHandler := handler.NewNodeHandler()

	// API路由组
	api := r.Group("/api")
	{
		// 认证接口 (无需token)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/send-register-code", authHandler.SendRegisterCode)
			auth.POST("/send-reset-code", authHandler.SendResetCode)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// 需要认证的接口
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户接口
			user := authorized.Group("/user")
			{
				user.GET("/info", userHandler.GetInfo)
				user.PUT("/info", userHandler.UpdateInfo)
				user.POST("/change-password", userHandler.ChangePassword)
			}

			// 账户接口
			account := authorized.Group("/account")
			{
				account.GET("/balance", userHandler.GetBalance)
				account.GET("/traffic", userHandler.GetTraffic)
				account.GET("/stats", userHandler.GetStats)
				account.POST("/recharge", subscriptionHandler.Recharge)
			}

			// 订阅接口
			subscriptions := authorized.Group("/subscriptions")
			{
				subscriptions.GET("", subscriptionHandler.GetList)
				subscriptions.GET("/plans", subscriptionHandler.GetPlans)
				subscriptions.POST("/purchase", subscriptionHandler.Purchase)
				subscriptions.POST("/renew", subscriptionHandler.Renew)
				subscriptions.DELETE("/:id", subscriptionHandler.Cancel)
			}

			// 节点接口
			nodes := authorized.Group("/nodes")
			{
				nodes.GET("", nodeHandler.GetList)
				nodes.GET("/:id", nodeHandler.GetDetail)
				nodes.POST("/:id/test", nodeHandler.TestLatency)
			}

			// TODO: 其他接口
			// - 公告 /announcements
			// - 下载 /downloads
			// - 统计 /stats
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}
