package config

import (
	"fmt"
	"log"
	"time"

	"backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 添加重试机制
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		DB, err = connectDB()
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 5) // 等待5秒后重试
	}

	if err != nil {
		log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(&models.ScanHistory{})
	if err != nil {
		log.Printf("Database migration failed: %v\n", err)
	}
}

func connectDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/security_scan?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
