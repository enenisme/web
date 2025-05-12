package api

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	// 可以添加依赖服务
}

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应结构
type LoginResponse struct {
	Token  string      `json:"token"`
	User   interface{} `json:"user"`
	Expire int64       `json:"expire"` // 过期时间戳（秒）
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 更新最后登录时间
	global.DB.Model(&user).UpdateColumn("last_login", time.Now())

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成令牌失败: " + err.Error(),
		})
		return
	}

	// 返回登录成功信息
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
		},
		Expire: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
		})
		return
	}

	// 获取用户信息
	var user models.User
	if err := global.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"role":     user.Role,
	})
}
