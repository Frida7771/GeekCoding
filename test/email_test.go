package test

import (
	"fmt"
	"net/smtp"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"

	"github.com/jordan-wright/email"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "<frida16571@gmail.com>"
	e.To = []string{"li.bo77771@gmail.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码是：<b>123456</b>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "frida16571@gmail.com", os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com"))

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("send email success")
}
