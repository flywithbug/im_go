package handle

import (
	"net/http"
	"im_go/im"
)

type RouterHandler func(resp http.ResponseWriter, req *http.Request)
var routers = map[string]RouterHandler{
	"/system": handleSystem, //系统状态
	"/summary":im.Summary,


	"/register": handleRegister, //注册
	"/login": handleLogin,
	"/logout":handleLogout,  //退出登录
	"/query": handleQuery,   //根据昵称查询用户列表




}

func RegisterRouters(prefix []string) {
	dup := make(map[string]bool)
	for _, p := range prefix {
		dup[p] = true
	}
	if len(dup) == 0 {
		dup[""] = true
	}
	for router, handler := range routers {
		for  k := range dup {
			doRouterRegister("/"+k+router, handler)
		}
	}
}
func doRouterRegister(router string, handler RouterHandler) {
	http.HandleFunc(router, handler)
}
