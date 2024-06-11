package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// genarate jwt Token
func GenerateJWTToken(openID string, expirationTime time.Time, secretKey []byte) (string, error) {
	// 创建一个新的 JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置 token 的声明（payload）
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = openID
	claims["exp"] = expirationTime.Unix()

	// 使用密钥进行签名，生成最终的 JWT token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("生成 JWT token 失败: %v", err)
	}

	return tokenString, nil
}
