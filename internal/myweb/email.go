package myweb

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

var _email, _server, _port, _code string

func SendEmail(to string, code string, sub string, body string) error {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	em := email.NewEmail()
	em.From = "星语<2024588807@qq.com>"
	em.To = []string{to}
	em.Subject = sub
	em.Text = []byte(body)
	fmt.Println(_email, _server, _port, _code)
	err := em.Send(_server+":"+_port, smtp.PlainAuth("", _email, _code, _server))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("send successfully!")
	return nil
}

/*func SendEmail(mailTo string, code string) error {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": _email,
		"pass": _code,
		"host": _server,
		"port": _port,
	}

	port, _ := strconv.Atoi(mailConn["port"])
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "xx官方"))                    // 添加别名
	m.SetHeader("To", mailTo)                                                         // 发送给用户(可以多个)
	m.SetHeader("Subject", "star whispers")                                           // 设置邮件主题
	m.SetBody("text/html", "Your code is"+code)                                       // 设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"]) // 设置邮件正文
	err := d.DialAndSend(m)
	return err
}*/

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
	_email = viper.GetString("email.email")
	_server = viper.GetString("email.server")
	_port = viper.GetString("email.port")
	_code = viper.GetString("email.auth_code")
}
