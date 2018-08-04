package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"im_go/model"
	"strings"
	"github.com/pborman/uuid"
	"im_go/common"
	"im_go/rsa"
	"im_go/config"
	log "github.com/flywithbug/log4go"
	"encoding/base64"
	"im_go/mail"
)

// 注册请求
/*
	Para:appId,account,password,Nick，头像地址。
*/



func handleRegister(c *gin.Context) {
	aRes := NewResponse()

	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	register := loginoutModel{}
	err := c.BindJSON(&register)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if len(register.VerifyKey) == 0 || len(register.Verify) == 0{
		aRes.SetErrorInfo(http.StatusBadRequest ,"Verify Code can not be nil")
		return
	}

	if !VerifyCaptcha(register.VerifyKey,register.Verify) {
		log.Info("%s %s",register.Verify,register.VerifyKey)
		aRes.SetErrorInfo(http.StatusBadRequest ,"Verify Code not right")
		return
	}

	if len(register.Account) < 4{
		aRes.SetErrorInfo(http.StatusBadRequest ,"account to short")
		return
	}
	if len(register.Password) < 6 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"password too shot")
		return
	}
	if register.Nick == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"nick can not be nil")
		return
	}
	if len(register.Mail) == 0{
		aRes.SetErrorInfo(http.StatusBadRequest,"Mail not right ")
		return
	}

	num ,err := model.CheckAccount(register.Account)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error"+err.Error())
		return
	}
	if num > 0{
		aRes.SetErrorInfo(http.StatusBadRequest,"Account already existed ")
		return
	}
	if register.AppId == 0{
		register.AppId = 10
	}
	password := common.Md5(register.Password)
	userId,err := model.SaveUser(register.AppId,register.Account,password,register.Password,register.Nick,register.Avatar,register.Mail)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError,"server error ")
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"register success")
	if len(register.Mail) != 0 && mail.MailStringVerify(register.Mail) {
		sendVerifyMail(register.Mail,*userId,register.Account,0)
	}
	//num ,_ = model.CheckAccount(register.Account)
	//if num > 0 {
	//	aRes.SetSuccessInfo(http.StatusOK,"register success")
	//	 return
	//}

	//aRes.SetErrorInfo(http.StatusInternalServerError,"server error")
}


/**
登录请求处理方法
*/
func handleLogin(c *gin.Context) {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()

	login := loginoutModel{}
	err := c.BindJSON(&login)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if login.Account == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"account can not be nil")
		return
	}
	if login.Password == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"password can not be nil")
		return
	}
	if login.Signature == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Signature can not be nil")
		return
	}
	//先decode key字符串
	decodeBytes, err := base64.StdEncoding.DecodeString(login.Signature)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"base64 DecodeString failure "+err.Error())
		return
	}
	//获取原始key
    b ,err:= rsa.RsaDecrypt(decodeBytes,config.Conf().RSAConfig.Private)
	if err != nil {
		log.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest ,"Signature verification failure ")
		return
	}
	signature := string(b)
	if  !strings.EqualFold(signature,login.Account+"-"+login.Password){
		aRes.SetErrorInfo(http.StatusBadRequest ,"Signature verification failure equal")
		return
	}

	num, err := model.CheckAccount(login.Account)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
		return
	}
	if num > 0{
		password := common.Md5(login.Password)
		user, err := model.LoginUser(login.Account, password)
		if err != nil {
			aRes.SetErrorInfo(http.StatusBadRequest ,"account or password not right")
			return
		}

		if !strings.EqualFold(user.UserId, "") {
			token := uuid.New()
			ip := common.GetIp(c.Request)
			userAgent ,ok := c.Get(KeyUserAgent)
			if !ok {
				userAgent = ""
			}
			deviceId ,_ := UserDeviceId(c)
			//log.Info(deviceId)
			if err := model.SaveLogin(user.GetAppId(),user.Uid,user.UserId, token, ip,user.Forbidden,userAgent.(string),deviceId); err != nil {
				log.Info(err.Error())
				aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
				return
			}
			user.Token = token
			aRes.SetSuccessInfo(http.StatusOK,"LoginSuccess")
			aRes.AddResponseInfo("user",user)
			return
		}
		log.Info("%s",user)
		aRes.SetErrorInfo(http.StatusInternalServerError ,"user not found")

	}else {
		aRes.SetErrorInfo(http.StatusBadRequest ,"account not existed"+err.Error())
	}
}


func handleLogout(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	logout := loginoutModel{}
	err := c.BindJSON(&logout)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if logout.Token == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"token can not be nil")
		return
	}
	_ ,err = model.Logout(logout.Token)
	if err != nil{
		errStr := err.Error()
		aRes.SetErrorInfo(http.StatusInternalServerError ,errStr)
		return
	}
	deviceId ,_ := UserDeviceId(c)
	model.DeviceUniteByUserId(deviceId)
	aRes.SetSuccessInfo(http.StatusOK,"success")
}




/**
查询请求处理方法
*/
func handleQueryNick(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	login := loginoutModel{}
	err := c.BindJSON(&login)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if login.Nick == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"nick can not be nil")
		return
	}
	users,err := model.QueryUserByNick(login.Nick)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
		return
	}
	aRes.AddResponseInfo("users",users)
}
/**
查询请求处理方法
*/
func handleQueryAccount(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	login := loginoutModel{}
	err := c.BindJSON(&login)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if login.Account == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"nick can not be nil")
		return
	}
	users,err := model.GetUserByAccount(login.Account)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
		return
	}
	aRes.AddResponseInfo("users",users)
}


func handleGetUserInfo(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	userId := c.Query("id")
	if len(userId) == 0{
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid")
		return
	}
	user,err := model.GetUserByUserId(userId)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
		return
	}
	aRes.AddResponseInfo("user",user)
}

func UpdateUserNickHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	para := loginoutModel{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}

	if para.Nick == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"nick can not be nil")
		return
	}
	user , _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"no user")
		return
	}
	err = model.UpdateUserNickName(para.Nick,user.UserId)
	if err != nil {
		log.Info(err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error"+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}

func ChangePasswordHandler(c *gin.Context) {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	para := loginoutModel{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}


	if len(para.OldPassword) == 0  {
		aRes.SetErrorInfo(http.StatusBadRequest, "oldpassword can not be nil")
		return
	}
	if len(para.NewPassword) < 6 {
		aRes.SetErrorInfo(http.StatusBadRequest, "passworld length less than 6")
		return
	}
	user, _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "no user")
		return
	}
	password := common.Md5(para.NewPassword)
	err = model.UpdateuserPassWordByOld(para.OldPassword, password, para.NewPassword, user.UserId)
	if err != nil {
		log.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "server error"+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}


func ChangePasswordByVerifyCodeHandler(c *gin.Context) {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	para := loginoutModel{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}
	if len(para.Verify) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest, "verify code can not be nil")
		return
	}
	if len(para.Account) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest, "Account can not be nill")
		return
	}
	if len(para.NewPassword) < 6 {
		aRes.SetErrorInfo(http.StatusBadRequest, "passworld length less than 6")
		return
	}
	_ ,err =  model.CheckVerifyByAccount(para.Account,para.Verify,1)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Verify invalid"+err.Error())
		return
	}
	password := common.Md5(para.NewPassword)
	err = model.UpdateuserPassWordByAccount(password, para.NewPassword, para.Account)
	if err != nil {
		log.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "server error"+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}


func UpdateUserCurrentLocation(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	para := loginoutModel{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}

	if len(para.Latitude) == 0 || len(para.Longitude) == 0{
		aRes.SetErrorInfo(http.StatusBadRequest, "Longitude Latitude can not be nil")
		return
	}
	user, _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "no user")
		return
	}
	err = model.UpdateUserLocations(para.Longitude,para.Latitude,para.LTimeStamp,user.UserId)
	if err != nil {
		log.Info(err.Error()+para.Latitude + " " +para.Longitude)
		aRes.SetErrorInfo(http.StatusBadRequest, "server error"+err.Error())
		return
	}
	if len(para.PIdentifier) == 0 {
		para.PIdentifier = uuid.NewUUID().String()
	}
	deviceId ,_ := UserDeviceId(c)
	err = model.SaveLocationsPath(user.UserId,para.Longitude,para.Latitude,para.LTimeStamp,para.PIdentifier,deviceId,para.LType)
	if err != nil {
		log.Info(err.Error())
	}

	aRes.SetSuccessInfo(http.StatusOK,"success")
}

func UpdateUserBatchLocations(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	para := PhotoLocationsModel{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}
	user, _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "no user")
		return
	}
	deviceId ,_ := UserDeviceId(c)
	for _,location := range para.List {
		if len(location.PIdentifier) > 0 {
			model.SaveLocationsPath(user.UserId,location.Longitude,location.Latitude,location.LTimeStamp,location.PIdentifier,deviceId,location.LType)
		}
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}