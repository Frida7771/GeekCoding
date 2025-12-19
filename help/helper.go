package help

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

// generate md5
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 生成 token
func GenerateToken(identity string, name string) (string, error) {

	var myKey = []byte("GeekCoding-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Identity: identity,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// 解析 token
func AnalyzeToken(tokenString string) (*UserClaims, error) {
	userClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte("GeekCoding-key"), nil
	})

	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyze token error: %v", err)
	}
	return userClaims, nil
}
