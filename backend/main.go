package main

import (
	"backend/config"
	"backend/routes"
	"log"

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

	// 启动服务器
	log.Println("Server starting on :8888...")
	if err := r.Run(":8888"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
