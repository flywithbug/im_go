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
		c.JSON(aRes.Code,aRes)
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
	_,err = model.SaveUser(login.AppId,login.Account,password,login.Nick,login.Avatar)
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
		c.JSON(aRes.Code,aRes)
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
	if login.Key == "" {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Key can not be nil")
		return
	}
	//先decode key字符串
	decodeBytes, err := base64.StdEncoding.DecodeString(login.Key)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Signature verification failure base64")
		return
	}
	//获取原始key
    b ,err:= rsa.RsaDecrypt(decodeBytes,config.Conf().RSAConfig.Private)
	if err != nil {
		log.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest ,"Signature verification failure r")
		return
	}
	key := string(b)
	if  !strings.EqualFold(key,login.Account+"-"+login.Password){
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
			aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
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
		aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
		return
	}else {
		aRes.SetErrorInfo(http.StatusBadRequest ,"account not existed"+err.Error())
	}
}


func handleLogout(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
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
		c.JSON(aRes.Code,aRes)
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

