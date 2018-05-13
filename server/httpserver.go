package server

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)


//TODO 路由中间件 auth校验

func StartHttpServer(address string,rPrefix []string)  {

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

	RegisterRouters(r,rPrefix)

	r.Run(address)
}



