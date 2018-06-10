package handler

import (
	"github.com/gin-gonic/gin"
	"im_go/config"
)


func AppConfigHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()
	aRes.AddResponseInfo("app_config",config.Conf().AppConfig)
}