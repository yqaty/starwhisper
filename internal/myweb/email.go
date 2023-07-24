package myweb

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

var _email, _server, _port, _code, _domain string

func SendEmail(to string, code string) error {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	em := email.NewEmail()
	em.From = "<" + _email + ">"
	em.To = []string{to}
	em.Subject = _domain + " Verification Code"
	em.Text = []byte("Your verification code is " + code + ", Please input within one minute.")
	err := em.Send(_server+":"+_port, smtp.PlainAuth("", _email, _code, _server))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("send successfully!")
	return nil
}

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
	_domain = viper.GetString("domain")
}
