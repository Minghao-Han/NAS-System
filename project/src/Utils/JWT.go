package Utils

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct {
	key string
}

var defaultJwt *JWT

func DefaultJWT() *JWT {
	if defaultJwt == nil {
		defaultJwt = &JWT{}
	}
	jwtKeyPath := DefaultConfigReader().Get("jwtKeyPath").(string)
	defaultJwt.key = YmlReader(jwtKeyPath, "key").(string)
	return defaultJwt
}

// 生成JWT Token
func (j JWT) GenerateToken(userId int) (string, error) {
	// Token 过期时间，这里设置为 24 小时
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 JWT 的 Claim
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime.Unix(),
	}

	// 使用指定的算法创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名 Token
	signedToken, err := token.SignedString([]byte(j.key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// 验证 JWT Token
func (j JWT) VerifyToken(tokenString string) (bool, error) {
	// 解析 Token，但不验证签名
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 返回所使用的密钥
		return []byte(j.key), nil
	})
	if err != nil {
		return false, err
	}

	// 验证 Token 是否有效
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func (j JWT) ExtractClaim(tokenString string, claimName string) (interface{}, error) {
	// 解析 Token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// 获取 Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	// 提取指定的 Claim 值
	if value, ok := claims[claimName]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("claim not found")
}

//
//func main() {
//	// 假设我们要生成并验证的用户名为 alice，密钥为 abc123
//	username := "alice"
//	secretKey := "abc123"
//
//	// 生成 Token 并打印出来
//	token, err := GenerateToken(username, secretKey)
//	if err != nil {
//		fmt.Printf("Failed to generate token: %v\n", err)
//		return
//	}
//	fmt.Printf("Generated token: %s\n", token)
//
//	// 验证 Token 是否有效
//	valid, err := verifyToken(token, secretKey)
//	if err != nil {
//		fmt.Printf("Failed to verify token: %v\n", err)
//		return
//	}
//	fmt.Printf("Token is valid: %t\n", valid)
//}
