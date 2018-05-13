package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"im_go/model"
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
	_,err = model.SaveUser(login.AppId,login.Account,login.Password,login.Nick,login.Avatar)
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