package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // 密码不会在JSON中返回
	Name      string    `gorm:"type:varchar(50)" json:"name"`        // 显示名称
	Email     string    `gorm:"type:varchar(100)" json:"email"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"` // 角色: admin, user
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// HashPassword 哈希用户密码
func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证用户密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
