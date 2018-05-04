package server

import (
	log "github.com/flywithbug/log4go"
	"fmt"
	"net/http"
	"im_go/server/handle"
)


//TODO 路由中间件 auth校验

func StartHttpServer(port int,routPrefix []string)  {
	log.Info("Http服务器启动中...")
	// 设置请求映射地址及对应处理方法
	handle.RegisterRouters(routPrefix)
	//打印监听端口
	log.Info("Http服务器开始监听[%d]端口", port)
	log.Info("*********************************************")
	// 设置监听地址及端口
	addr := fmt.Sprintf("localhost:%d", port)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(fmt.Errorf("监听Http失败: %s", err))
		}
	}()
}



