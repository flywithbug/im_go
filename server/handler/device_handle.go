package handler

import (
	"github.com/gin-gonic/gin"
	"im_go/model"
	"net/http"
	"github.com/flywithbug/log4go"
)

func RegistPushService(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()

	device := model.Device{}
	err := c.BindJSON(&device)
	if err != nil {
		log4go.Info(err.Error())
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
	err = device.SaveOrUpdateToDb()
	if err != nil {
		errStr := err.Error()
		aRes.SetErrorInfo(http.StatusInternalServerError ,errStr)
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}


func GetPushStatusHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()

	deviceId := c.Param("id")
	if len(deviceId) < 10 {
		aRes.SetErrorInfo(http.StatusBadRequest ,"deviceId invalid")
		return
	}
	d, err := model.GetDevicesByDeviceId(deviceId)
	if err != nil {
		errStr := err.Error()
		aRes.SetErrorInfo(http.StatusNotFound ,errStr)
		return
	}
	aRes.AddResponseInfo("device",d)

}

