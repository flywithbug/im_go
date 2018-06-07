package server

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"fmt"
	"im_go/server/handler"
)


//TODO 路由中间件 auth校验

func StartHttpServer(address string,rPrefix []string)  {

	//log.Info("Http服务器启动中...")
	//// 设置请求映射地址及对应处理方法
	////打印监听端口
	//log.Info("Http服务器开始监听[%s]端口", address)
	//log.Info("*********************************************")
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	r.Use(handler.TokenAuthMiddleware())
	handler.RegisterRouters(r,rPrefix)
	err := r.Run(address)
	if  err != nil {
		panic(fmt.Errorf("监听Http失败: %s", err))
	}

}



