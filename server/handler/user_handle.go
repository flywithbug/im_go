package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"im_go/model"
	"strings"
	"github.com/pborman/uuid"
	"im_go/common"
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
			if err := model.SaveLogin(user.GetAppId(),user.Uid,user.UserId, token, ip,user.Forbidden); err != nil {
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
	num ,err := model.Logout(login.Token)
	if num <= 0 || err != nil{
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


/*
	三个参数，
	u_id 好友请求发送者
	friend_id 接收人
	method  "add",delete,"remark" 添加备注
*/
func handleAddRelation(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()
	m := relationShipModel{}
	err := c.BindJSON(&m)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if strings.EqualFold(m.Method,"add") {
		_, err = model.AddUserRelation(m.UId,m.FriendId)
		if err != nil {
			aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
			return
		}else {
			aRes.SetSuccessInfo(http.StatusOK,"success")
			return
		}
	}else if strings.EqualFold(m.Method,"delete"){
		err = model.DelRelationShip(m.RelationId)
		if err != nil {
			aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
			return
		}else {
			aRes.SetSuccessInfo(http.StatusOK,"success")
			return
		}
	}else if strings.EqualFold(m.Method,"remark") { //删除好友
		err = model.UpdateRelationRemark(m.RelationId,m.Remark)
		if err != nil {
			aRes.SetErrorInfo(http.StatusInternalServerError ,"server error")
			return
		}else {
			aRes.SetSuccessInfo(http.StatusOK,"success")
			return
		}
	}else {
		aRes.SetSuccessInfo(http.StatusNotFound,"Not Found: "+m.Method)
	}
}
