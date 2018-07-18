package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)



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
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configC)
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	return idKeyD,base64stringD
}

func VerfiyCaptcha(verifyKey,verifyValue string) bool{
	return base64Captcha.VerifyCaptcha(verifyKey, verifyValue)
}
