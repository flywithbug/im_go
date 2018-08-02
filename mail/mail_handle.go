package mail

import (
	"gopkg.in/gomail.v2"
	"im_go/config"
	"errors"
)

var(
	Mail *gomail.Dialer
)

func SendMail(to ,title ,from string,mType int,body string) error  {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", config.Conf().MailConfig.Username,title)
	m.SetHeader("Subject", from)
	m.SetHeader("To", to)
	m.SetBody("text/html", body)

	if Mail == nil{
		return errors.New("mail is nil")
	}
	return Mail.DialAndSend(m)
}

func SendVerifyMail(uuid,mail string)  {
	//verifyStr := "www.flywithme.top:"
}

