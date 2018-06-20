package im

import (
	"im_go/model"
	"fmt"
)

/*
	此服务需要最后启动，
*/
func StartIMServer(im_port int,http_port int)  {
	if model.Database == nil {
		panic(fmt.Errorf("mysql服务未连接"))
	}
	//启动网络无法
	startHttpServer(http_port)
	//启动IM服务
	listen(im_port)
}