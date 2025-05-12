package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// 密钥，实际应用中应该从环境变量或配置文件读取
	jwtSecret = []byte("web_vulnerability_scanner_secret_key")
)

// Claims 自定义JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username, role string) (string, error) {
	// 设置令牌有效期为24小时
	expireTime := time.Now().Add(24 * time.Hour)

	// 创建JWT声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_token",
		},
	}

	// 创建令牌
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(token string) (*Claims, error) {
	// 解析令牌
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 检查令牌是否有效
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}
