package handler

import (
	"github.com/gin-gonic/gin"
	"strings"
	"im_go/model"
	"net/http"
	"im_go/config"
)

//todo  添加签名进行客户端请求安全校验
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlPath := ctx.Request.URL.Path
		if strings.EqualFold(urlPath,KeyLoginPath){
			userAg := ctx.GetHeader(KeyUserAgent)
			if len(userAg)> 0 {
				ctx.Set(KeyUserAgent,userAg)
			}
			return
		}

		for _,v := range config.Conf().AuthFilterWhite {
			for _, pr := range config.Conf().RouterPrefix{
				if strings.HasPrefix(urlPath,pr + v) {
					return
				}
			}
		}
		//log4go.Info(urlPath)
		aResp := NewResponse()
		token,_ := ctx.Cookie(KeyUserToken)
		//log4go.Info("TokenAuthMiddleware",token,ctx.Request.Header)
		var aUser *model.User
		if len(token) > 0 {
			aUser, _ = model.GetUserByToken(token)
		}else {
			aResp.SetErrorInfo(http.StatusUnauthorized,"no token found")
			ctx.JSON(http.StatusOK,aResp)
			ctx.Abort()
			return
		}
		if aUser == nil {
			aResp.SetErrorInfo(http.StatusUnauthorized,"auth token invalid")
			ctx.JSON(http.StatusOK,aResp)
			ctx.Abort()
			return
		}
		ctx.Set(KeyContextUser,aUser)
	}
}