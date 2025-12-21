package help

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"net/smtp"
	"os"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.RegisteredClaims
}

// generate md5
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 生成 token
func GenerateToken(identity string, name string, isAdmin int) (string, error) {

	var myKey = []byte("GeekCoding-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Identity: identity,
		Name:     name,
		IsAdmin:  isAdmin,
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

// 发送邮件
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "<frida16571@gmail.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码是：<b>" + code + "</b>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "frida16571@gmail.com", os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com"))
	if err != nil {
		return err
	}
	fmt.Println("send code success")
	return nil
}

// get uuid
func GetUUID() string {
	return uuid.NewV4().String()
}

// generate random code
func GetRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// save code
func SaveCode(code []byte) (string, error) {
	// 使用系统临时目录，避免影响 swag init
	tmpDir := os.TempDir()
	dirName := tmpDir + "/code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.MkdirAll(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil

}
