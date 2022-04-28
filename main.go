package main

import (
	"github.com/xinliangnote/go-util/mail"
	"os"
	"sign/conf"
	"sign/iqiyi"
	"time"
)

func init() {
	conf.Setup()
}

func main() {
	iqiyiCookie := os.Getenv("iqiyiCookie")
	for {
		err := iqiyi.ParseCookie(iqiyiCookie).DoSomeThings()
		if err != nil {
			sendMail(err.Error())
		}
		time.Sleep(time.Hour * 24)
	}
}

func sendMail(content string) {
	subject := "自动签到签到失败"
	body := content + "，请尽快处理！"
	if conf.Email.Host != "" {
		options := &mail.Options{
			MailHost: conf.Email.Host,
			MailPort: conf.Email.Port,
			MailUser: conf.Email.User,
			MailPass: conf.Email.Pass,
			MailTo:   conf.Email.AdminUser,
			Subject:  subject,
			Body:     body,
		}
		_ = mail.Send(options)
	}
}

