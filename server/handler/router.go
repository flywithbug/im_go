package handler

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"im_go/im"
	"strings"
)

const  LoginPATH  =  "/login"



var routers = map[string]gin.HandlerFunc{
	"GET   		/system": 				handleSystem, //系统状态
	"GET		/summary":				im.GinSummary,
	"GET		/config":				AppConfigHandler, //请求网络接口下发，（最好使用使用域名请求）
	"POST 		/register": 			handleRegister, //注册
	"POST  		/login": 				handleLogin,
	"POST		/logout":				handleLogout,  //退出登录
	"POST 		/query": 				handleQuery,   //根据昵称查询用户列表
	"GET 		/user/:id": 			handleGetUserInfo,   //根据用户user_id获取用户信息
	"POST		/upload/avatar":		UploadImageHandler,			//avatar图片上传
	"GET	    /image/avatar/":		DownloadImageHandler,      //avatar图片下载
	"POST		/push":					RegistPushService,
	"POST       /push/authorization": 	UpdatePushStatusService,  //更新推送权限

}

func RegisterRouters(r *gin.Engine, prefixs []string){
	dup := make(map[string]bool)
	for _, p := range prefixs {
		dup[p] = true
	}
	if len(dup) == 0 {
		dup[""] = true
	}
	for router, handler := range routers {
		method ,path := regexpRouters(router)
		for  k := range dup {
			funcDoRouteRegister(method,strings.ToLower(k+path),handler,r)//path 全小写
		}
	}
}

func funcDoRouteRegister(method, route string, handler gin.HandlerFunc, r *gin.Engine) {
	switch method {
	case "POST":
		r.POST(route, handler)
	case "PUT":
		r.PUT(route, handler)
	case "HEAD":
		r.HEAD(route, handler)
	case "DELETE":
		r.DELETE(route, handler)
	case "OPTIONS":
		r.OPTIONS(route, handler)
	default:
		r.GET(route, handler)
	}
}

var routerRe = regexp.MustCompile(`([A-Z]+)[^/]*(/[a-zA-Z/:]+)`)
func regexpRouters(router string) (method,path string) {
	match := routerRe.FindAllStringSubmatch(router, -1)
	return match[0][1],match[0][2]
}
