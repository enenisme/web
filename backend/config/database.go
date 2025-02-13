package config

import (
	"fmt"
	"log"

	"backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:root@tcp(127.0.0.1:3306)/security_scan?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(&models.ScanHistory{})
	if err != nil {
		fmt.Printf("Database migration failed: %v\n", err)
	}
}
