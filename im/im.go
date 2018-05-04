package im



/*
	此服务需要最后启动，
*/
func StartIMServer(im_port int,http_port int)  {
	//启动网络无法
	startHttpServer(http_port)


	//启动IM服务
	listen(im_port)
}