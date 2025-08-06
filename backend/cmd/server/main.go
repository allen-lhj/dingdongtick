package main

import (
	"log"
	"net/http"

	"ticktick-backend/config"
	"ticktick-backend/internal/dal"
	"ticktick-backend/internal/handlers"
	"ticktick-backend/internal/middleware"
	"ticktick-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	db, err := dal.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 自动迁移数据库
	// if err := db.AutoMigrate(); err != nil {
	// 	log.Fatalf("数据库迁移失败: %v", err)
	// }

	// 初始化Redis服务
	redisService, err := services.NewRedisService(&cfg.Redis)
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	defer redisService.Close()

	// 初始化Token存储服务
	tokenStore := services.NewTokenStore(redisService)

	// 初始化Token监控服务
	tokenMonitor := services.NewTokenMonitor(tokenStore, redisService)
	tokenMonitor.Start()
	defer tokenMonitor.Stop()

	// 初始化服务层
	userService := services.NewUserService(db)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(userService, tokenStore, cfg)
	monitorHandler := handlers.NewMonitorHandler(tokenMonitor, tokenStore)

	// 创建Gin路由器
	router := gin.Default()

	// 配置CORS中间件 - 允许所有端口
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5273", "http://127.0.0.1:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// API路由组
	api := router.Group("/api/v1")

	// 健康检查端点
	api.GET("/health", monitorHandler.GetHealth)
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// 认证路由（不需要JWT验证）
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", authHandler.Logout)

	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg, tokenStore))
	{
		// 用户信息路由
		protected.GET("/profile", authHandler.GetProfile)

		// 会话管理路由
		protected.GET("/sessions", authHandler.GetSessions)
		protected.DELETE("/sessions/:tokenId", authHandler.RevokeSession)
		protected.POST("/logout-all", authHandler.LogoutAll)

		// 管理员功能（实际项目中需要权限验证）
		protected.POST("/revoke-token", authHandler.RevokeToken)

		// 监控相关路由
		monitor := protected.Group("/monitor")
		{
			monitor.GET("/health", monitorHandler.GetHealth)
			monitor.GET("/stats", monitorHandler.GetStats)
			monitor.GET("/data", monitorHandler.GetMonitoringData)
			monitor.POST("/cleanup", monitorHandler.ForceCleanup)
			monitor.GET("/token/:tokenId", monitorHandler.GetTokenInfo)
			monitor.GET("/user/:userId/sessions", monitorHandler.GetUserSessions)
			monitor.DELETE("/user/:userId/sessions", monitorHandler.RevokeUserAllSessions)
			monitor.GET("/metrics", monitorHandler.GetSystemMetrics)
		}

		// 项目路由
		projects := protected.Group("/projects")
		{
			projects.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "获取项目列表功能待实现"})
			})
			projects.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "创建项目功能待实现"})
			})
			projects.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "获取项目详情功能待实现"})
			})
			projects.PUT("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "更新项目功能待实现"})
			})
			projects.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "删除项目功能待实现"})
			})
		}

		// 任务路由
		tasks := protected.Group("/tasks")
		{
			tasks.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "获取任务详情功能待实现"})
			})
			tasks.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "创建任务功能待实现"})
			})
			tasks.PUT("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "更新任务功能待实现"})
			})
			tasks.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "删除任务功能待实现"})
			})
			tasks.POST("/:id/complete", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "完成任务功能待实现"})
			})
		}

		// 日历视图路由
		calendar := protected.Group("/calendar")
		{
			calendar.GET("/view", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "日历视图功能待实现"})
			})
		}
	}

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("服务器启动在 %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
