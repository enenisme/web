package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	// ... other imports
)

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 创建 Gin 引擎实例
	r := gin.Default()

	// 配置路由
	routes.SetupRoutes(r)

	// 优雅关闭
	go func() {
		// 启动服务器
		log.Println("Server starting on :8888...")
		if err := r.Run(":8888"); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
