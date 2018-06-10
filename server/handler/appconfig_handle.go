package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppConfigModel struct {
	AppId 			int			`json:"app_id"`
	Platform		int			`json:"platform"`   //1.iOS, 2.android
}



func AppConfigHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()
	app := AppConfigModel{}
	err := c.BindJSON(&app)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"Param invalid"+err.Error())
		return
	}


}