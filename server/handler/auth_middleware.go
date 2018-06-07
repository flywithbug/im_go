package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/flywithbug/log4go"
	"strings"
)


func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlPath := ctx.Request.URL.Path
		if strings.EqualFold(urlPath,KeyLoginPath){
			userAg := ""
			if c, ok := ctx.Request.Header[KeyUserAgent]; ok && len(c) > 0{
				userAg = strings.Join(c, "; ")
			}
			if len(userAg)> 0 {
				ctx.Set(KeyUserAgent,userAg)
			}
			return
		}
		cookieStr := ""
		if c, ok := ctx.Request.Header["Cookie"]; ok && len(c) > 0{
			cookieStr = strings.Join(c, "; ")
		}

		println("====> Path:", urlPath, "Cookie:", cookieStr)

		log.Info("%s",ctx.Request.Header)


		//ctx.Next()
	}

}