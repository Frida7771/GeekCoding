package test

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

// 生成 token
func TestGenerateToken(t *testing.T) {
	UserClaims := &UserClaims{
		Identity:         "user_1",
		Name:             "Get",
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	var myKey = []byte("GeekCoding-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)
	tokenString, err := token.SignedString(myKey)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenString)
}

// 解析 token
func TestAnalyzeToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.BBiU8x7h3LRsMOj53a8jHZds-18lUoV2EQa3pOURIXE"
	var myKey = []byte("GeekCoding-key")
	userClaims := &UserClaims{}
	claims, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})

	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userClaims)
	}
}
