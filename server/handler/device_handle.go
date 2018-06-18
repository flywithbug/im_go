package handler

import (
	"github.com/gin-gonic/gin"
	"im_go/model"
	"net/http"
)

func RegistPushService(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()

	device := model.Device{}
	err := c.BindJSON(&device)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}
	if len(device.DeviceToken) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"DeviceToken can not be nil")
		return
	}
	if len(device.DeviceId) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"DeviceId can not be nil")
		return
	}
	if len(device.UserAgent) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"UserAgent can not be nil")
		return
	}
	if device.Platform == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Platform must > 0,(iphone android,web),1/2/3")
		return
	}
	user ,_:= User(c)
	if user != nil {
		device.UserId = user.UserId
	}
	err = device.SaveToDb()
	if err != nil {
		errStr := err.Error()
		aRes.SetErrorInfo(http.StatusInternalServerError ,errStr)
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}