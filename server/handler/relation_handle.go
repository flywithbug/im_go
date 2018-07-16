package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"im_go/model"
	"strings"
)

func ApplyUserRelation(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	reModel := UserRelationShip{}
	err := c.BindJSON(&reModel)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if len(reModel.UserId) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"userId can not be nil")
		return
	}

	user , _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"no user")
		return
	}

	if strings.EqualFold(user.UserId,reModel.UserId){
		aRes.SetErrorInfo(http.StatusBadRequest ,"can not add self as friend")
		return
	}


	err = model.ApplyUserRelationShip(user.UserId,reModel.UserId,reModel.UserId)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}

func UpdateUserRelation(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	reModel := UserRelationShip{}
	err := c.BindJSON(&reModel)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if len(reModel.UserId) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"userId can not be nil")
		return
	}
	user , _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"no user")
		return
	}

	if strings.EqualFold(user.UserId,reModel.UserId){
		aRes.SetErrorInfo(http.StatusBadRequest ,"error param")
		return
	}

	err = model.UpdateUserRelation(user.UserId,reModel.UserId,reModel.Status)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}

func FriendsListHandle(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	user , _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"no user")
		return
	}
	list,err := model.GetFriendsList(user.UserId)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError ,err.Error())
		return
	}
	aRes.AddResponseInfo("list",list)
}