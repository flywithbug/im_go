package handle

import "net/http"

type RouterHandler func(resp http.ResponseWriter, req *http.Request)
var routers = map[string]RouterHandler{
	"/system": handleSystem, //系统状态

	"/register": handleRegister, //注册
	"/login": handleLogin,

	//"/query": handleQuery,
	//"/users/relation/add": handleUserRelationAdd,
	//"/users/relation/del": handleUserRelationDel,
	//"/users/relation/push": handleUserRelationPush,
	//"/users/relation/refuse": handleUserRelationRefuse,
	//"/users/category/add": handleUserCategoryAdd,
	//"/users/category/del": handleUserCategoryDel,
	//"/users/category/edit": handleUserCategoryEdit,
	//"/users/category/query": handleUserCategoryQuery,

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
