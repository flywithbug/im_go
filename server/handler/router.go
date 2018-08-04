package handler

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"im_go/im"
	"strings"
)

const  LoginPATH  =  "/login"



var routers = map[string]gin.HandlerFunc{
	"GET   		/system": 						handleSystem, //系统状态
	"GET		/summary":						im.GinSummary,
	"GET		/config":						AppConfigHandler, //请求网络接口下发，（最好使用使用域名请求）
	"POST 		/register": 					handleRegister, //注册
	"GET        /verify/generate": 				GenerateCaptchaHandler, //生成验证码
	"POST       /verify/mail":                  SendVerifyMailHandle,//发送邮箱校验
	"GET 		/mail/check":    		        VerifyMailHandle, //邮箱校验
	"GET        /account/mail":                 GetMailByAccountHandle,//获取验证邮箱

	"POST  		/login": 						handleLogin,
	"POST		/logout":						handleLogout,  //退出登录
	"POST 		/query/nick": 					handleQueryNick,   //根据昵称查询用户列表
	"POST 		/query/account": 				handleQueryAccount,   //根据账号查询用户
	"POST		/update/nick":					UpdateUserNickHandler, //更新用户昵称
	"GET 		/user/info": 					handleGetUserInfo,   //根据用户user_id获取用户信息
	"POST		/upload/avatar":				UploadUserAvatarHandler,			//avatar图片上传
	"GET	    /image/avatar/":				DownloadImageHandler,      //avatar图片下载
	"POST		/push/update":					RegistPushService,
	"GET    	/push/status/:id":				GetPushStatusHandler,    //当前推送状态获取 (根据deviceId 获取推送状态设置)
	"POST       /user/password":				ChangePasswordHandler,
	"POST 		/user/location":        		UpdateUserCurrentLocation, //更新当前位置信息
	"POST       /user/location/batch":  		UpdateUserBatchLocations,//批量更新定位信息
	"POST       /user/location/current":        GetUserCurrentLocations,

	"POST       /user/relation/apply":  		ApplyUserRelation, //申请好友
	"POST       /user/relation/update": 		UpdateUserRelation, //更新好友关系
	"GET		/user/friend/list":				FriendsListHandle,  //好友列表
	"POST       /user/authorization/update":	UpdateAuthorization, //更新用户权限
	"POST       /user/authorization/status":    GetAuthorizationStatus,//获取授权状态
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
	case "GET":
		r.GET(route, handler)
	default:
		r.GET(route, handler)
	}
}

var routerRe = regexp.MustCompile(`([A-Z]+)[^/]*(/[a-zA-Z/:]+)`)
func regexpRouters(router string) (method,path string) {
	match := routerRe.FindAllStringSubmatch(router, -1)
	return match[0][1],match[0][2]
}
