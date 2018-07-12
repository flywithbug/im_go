package handler

import (
	"github.com/gin-gonic/gin"
	"im_go/config"
	"net/http"
)


func AppConfigHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()
	aRes.AddResponseInfo("app_config",config.Conf().AppConfig)
}