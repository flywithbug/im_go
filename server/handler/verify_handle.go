package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"github.com/pborman/uuid"
	"im_go/model"
	log "github.com/flywithbug/log4go"
	"im_go/mail"
	"strings"
	"time"
)

type VerifyModel struct {
	Mail     string  `json:"mail"`
	UserId   string	 `json:"user_id"`
	UUID     string  `json:"uuid"`
	VType    int     `json:"v_type"`
	Account  string  `json:"account"`
}

func GenerateCaptchaHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	verifyKey,base64Png := codeCaptchaCreate()
	aRes.AddResponseInfo("base64",base64Png)
	aRes.AddResponseInfo("verifyKey",verifyKey)
}

func codeCaptchaCreate()(verifyKey, base64Png string)  {
	//var configD = base64Captcha.ConfigDigit{
	//	Height:     40,
	//	Width:      120,
	//	MaxSkew:    0,
	//	DotCount:   0,
	//	CaptchaLen: 4,
	//}
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height:             120,
		Width:              240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}

	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	uuIDS := uuid.NewUUID().String()
	idKeyD, capD := base64Captcha.GenerateCaptcha(uuIDS, configC)
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	return idKeyD,base64stringD
}

func VerifyCaptcha(verifyKey,verifyValue string) bool{
	return base64Captcha.VerifyCaptcha(verifyKey, verifyValue)
}

func sendVerifyMail(Mail, userId, account string, vType int) error {
	var vld int
	if vType == 1 {
		vld = int(time.Now().Unix() + 60*10)
	}
	uuId,vCode, err := model.GeneryVerifyData(userId,account,vld,vType)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	if vType == 1 {
		return mail.SendVerifyCode(vCode,Mail)
	}
	return mail.SendVerifyMail(uuId,Mail)
}


func SendVerifyMailHandle(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	verify := VerifyModel{}
	err := c.BindJSON(&verify)
	if err != nil {
		log.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if len(verify.Account) != 0 {
		user,err := model.GetMailByAccount(verify.Account)
		if err != nil {
			aRes.SetErrorInfo(http.StatusBadRequest ,"no user found "+err.Error())
			return
		}
		verify.Mail = user.Mail
		verify.UserId = user.UserId
	}

	if len(verify.Mail) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"mail invalid")
		return
	}

	if len(verify.UserId) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"UserId invalid")
		return
	}
	err = sendVerifyMail(verify.Mail,verify.UserId,verify.Account,verify.VType)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"mail server error "+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}


func VerifyMailHandle(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	uuId := c.Query("uuid")
	vType := c.Query("type")
	if len(uuId) < 10 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"uuid invalid")
		return
	}
	userId, err := model.CheckVerify(uuId,vType)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"no user found"+err.Error())
		return
	}
	err = model.UpdateUserMailVerifyChecked(userId)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"no user found"+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}

func GetMailByAccountHandle(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	account := c.Query("account")
	if len(account) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"account invalid")
		return
	}
	user ,err := model.GetMailByAccount(account)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
		return
	}
	mails := strings.Split(user.Mail,"@")
	if len(mails) == 2 {
		aRes.SetResponseDataInfo("mail",Substr(mails[0],0,3) + "*****" + "@" + mails[1])
	}else {
		aRes.SetErrorInfo(http.StatusBadRequest ,"mail invalid")
		return
	}
}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
