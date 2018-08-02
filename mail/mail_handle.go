package mail

import (
	"gopkg.in/gomail.v2"
	"im_go/config"
	"errors"
	"regexp"
	"strings"
)

var(
	Mail *gomail.Dialer
)

func sendMail(to ,title ,subject string,mType int,body string) error  {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", config.Conf().MailConfig.Username,title)
	m.SetHeader("Subject", subject)
	m.SetHeader("To", to)
	m.SetBody("text/html", body)

	if Mail == nil{
		return errors.New("mail is nil")
	}
	return Mail.DialAndSend(m)
}

func SendVerifyMail(uuid,mail string) error {
	verifyStr := "http://www.flywithme.top/im?uuid=" + uuid
	return sendMail(mail,"足迹","邮箱验证",0,verifyStr)
}

var routerRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`)
func MailStringVerify(mail string)bool  {
	match := routerRe.FindString(mail)
	return strings.EqualFold(match,mail)
}