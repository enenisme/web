package config

import (
	"fmt"
	"log"
	"time"

	"backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 添加重试机制
	var err error
	var db *gorm.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = connectDB()
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
	err = db.AutoMigrate(&models.ScanHistory{}, &models.User{})
	if err != nil {
		log.Printf("Database migration failed: %v\n", err)
	}

	// 初始化管理员用户
	initAdminUser(db)

	return db
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

// 初始化管理员用户
func initAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		// 创建默认管理员账户
		admin := models.User{
			Username: "admin",
			Password: "123456", // 初始密码
			Name:     "管理员",
			Email:    "admin@example.com",
			Role:     "admin",
		}

		// 哈希密码
		err := admin.HashPassword()
		if err != nil {
			log.Printf("Failed to hash admin password: %v\n", err)
			return
		}

		// 保存到数据库
		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("Failed to create admin user: %v\n", result.Error)
		} else {
			log.Println("Admin user created successfully")
		}
	}
}
