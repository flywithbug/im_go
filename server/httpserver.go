package server

import (
	"im_go/config"
	"log"
	"fmt"
	"net/http"
	"im_go/server/handler"
)

func StartHttpServer(conf config.IMConfig)error  {
	log.Printf("Http服务器启动中...")
	// 设置请求映射地址及对应处理方法
	handler.RegisterRouters(conf.RouterPrefix)
	//打印监听端口
	log.Printf("Http服务器开始监听[%d]端口", conf.HttpPort)
	log.Println("*********************************************")
	// 设置监听地址及端口
	addr := fmt.Sprintf("0.0.0.0:%d", conf.HttpPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("监听Http失败: %s", err)
	}
	return nil
}



