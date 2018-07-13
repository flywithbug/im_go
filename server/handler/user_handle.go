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
	if login.Nick == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"nick can not be nil")
		return
	}
	if len(login.Password) < 6 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"password too shot")
		return
	}
	num ,err := model.CheckAccount(login.Account)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error"+err.Error())
		return
	}
	if num > 0{
		aRes.SetErrorInfo(http.StatusBadRequest,"Account already existed ")
		return
	}
	if login.AppId == 0{
		login.AppId = 10
	}

	password := common.Md5(login.Password)
	_,err = model.SaveUser(login.AppId,login.Account,password,login.Password,login.Nick,login.Avatar)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError,"server error ")
		return
	}
	num ,_ = model.CheckAccount(login.Account)
	if num > 0 {
		aRes.SetSuccessInfo(http.StatusOK,"register success")
		 return
	}
	aRes.SetErrorInfo(http.StatusInternalServerError,"server error")
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
			if err := model.SaveLogin(user.GetAppId(),user.Uid,user.UserId, token, ip,user.Forbidden,userAgent.(string)); err != nil {
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
	login := loginoutModel{}
	err := c.BindJSON(&login)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if login.Token == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"token can not be nil")
		return
	}
	_ ,err = model.Logout(login.Token)
	if err != nil{
		errStr := err.Error()
		aRes.SetErrorInfo(http.StatusInternalServerError ,errStr)
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")

}




/**
查询请求处理方法
*/
func handleQuery(c *gin.Context)  {
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
	users,err := model.QueryUser(login.Nick)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
		return
	}
	aRes.AddResponseInfo("users",users)
}

func handleGetUserInfo(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()
	userId := c.Param("id")
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
		aRes.SetErrorInfo(http.StatusBadRequest ,err.Error())
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
	if len(para.OldPassword) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest, "oldpassword can not be nil")
		return
	}
	if len(para.NewPassword) < 6 {
		aRes.SetErrorInfo(http.StatusBadRequest, "passworld length less than 6")
		return
	}
	user, _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	password := common.Md5(para.NewPassword)
	err = model.UpdateuserPassWord(para.OldPassword, password, para.NewPassword, user.UserId)
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
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	err = model.UpdateUserLocations(para.Longitude,para.Latitude,para.LTimeStamp,user.UserId)
	if err != nil {
		log.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "server error"+err.Error())
		return
	}
	err = model.SaveLocationsPath(user.UserId,para.Longitude,para.Latitude,para.LTimeStamp,para.LType)
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
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}

	for _,location := range para.List {
		err = model.SaveLocationsPath(user.UserId,location.Longitude,location.Latitude,location.LTimeStamp,location.LType)
		if err != nil {
			log.Info(err.Error())
			break
		}
	}

	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError, err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")

}