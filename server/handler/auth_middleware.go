package handler

import (
	"github.com/gin-gonic/gin"
	"strings"
	"im_go/model"
	"net/http"
	"im_go/config"
	"fmt"
)

var whithList []string

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
			fmt.Println(v,urlPath)
			if strings.HasSuffix(urlPath,v) {
				return
			}
		}

		aResp := NewResponse()
		token,_ := ctx.Cookie(KeyUserToken)
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

func WhitePathList()[]string  {
	if whithList == nil{
		whithList = make([]string,0,len(config.Conf().AuthFilterWhite)*len(config.Conf().RouterPrefix))

		for i,v := range config.Conf().AuthFilterWhite {
			for _,vk := range config.Conf().RouterPrefix {

				whithList[i]= fmt.Sprintf("/%s%s",vk,v)
			}
		}
	}
	return whithList
}
