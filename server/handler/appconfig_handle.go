package handler

import (
	"github.com/gin-gonic/gin"
	"im_go/config"
)

//type AppConfigModel struct {
//	AppId 			int			`json:"app_id"`
//	Platform		int			`json:"platform"`   //1.iOS, 2.android
//}


//type AppConfigResponse struct {
//	ApiHost 		string  	`json:"api_host"`    //api请求host
//	IMSocketHost	string		`json:"im_socket_host"` //IM通讯host
//	IMSocketPort    int			`json:"im_socket_port"` //IM通讯 port
//	DomainName		string		`json:"domain_name"`  //域名
//}


func AppConfigHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()
	aRes.AddResponseInfo("app_config",config.Conf().AppConfig)
}