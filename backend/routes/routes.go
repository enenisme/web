package routes

import (
	"backend/api"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8082"}, // 允许前端域名访问
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 初始化所有 handlers
	fingerprintHandler := api.NewFingerprintHandler()
	hostAliveHandler := api.NewHostAliveHandler()
	portScanHandler := api.NewPortScanHandler()
	scanHistoryHandler := api.NewScanHistoryHandler()
	subdomainHandler := api.NewSubdomainHandler()
	vulnerabilityHandler := api.NewVulnerabilityHandler()

	// 基础测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 指纹识别相关路由
		apiGroup.POST("/fingerprint/scan", fingerprintHandler.HandleFingerprintScan)

		// 主机存活检测相关路由
		apiGroup.POST("/host/alive", hostAliveHandler.HandleHostAliveCheck)

		// 端口扫描相关路由
		apiGroup.POST("/port/scan", portScanHandler.HandlePortScan)

		// 子域名扫描相关路由
		apiGroup.POST("/subdomain/scan", subdomainHandler.HandleSubdomainScan)

		// 漏洞扫描相关路由
		apiGroup.POST("/vulnerability/scan", vulnerabilityHandler.HandleVulnerabilityScan)

		// 扫描历史相关路由
		historyGroup := apiGroup.Group("/history")
		{
			historyGroup.GET("", scanHistoryHandler.GetScanHistoryList)       // 获取历史列表
			historyGroup.GET("/:id", scanHistoryHandler.GetScanHistoryDetail) // 获取历史详情
			historyGroup.DELETE("/:id", scanHistoryHandler.DeleteScanHistory) // 删除历史记录
		}
	}
}
