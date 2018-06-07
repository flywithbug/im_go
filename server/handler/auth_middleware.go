package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/flywithbug/log4go"
	"strings"
	"fmt"
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
		token := CookieToMap(cookieStr)[KeyUserToken]
		if len(token) > 0 {

		}



		println("====> Path:", urlPath, "Cookie:", cookieStr,)
		fmt.Println(CookieToMap(cookieStr))
		log.Info("%s ",ctx.Request.Header)

		//ctx.Next()
	}

}

func CookieToMap(Cookie string) map[string]string {
	keyValues := strings.Split(Cookie, ";")
	m := make(map[string]string)
	for _,v:= range keyValues{
		kvs := strings.Split(v,"=")
		if len(kvs) == 2{
			m[kvs[0]]=kvs[1]
		}
	}
	return m
}