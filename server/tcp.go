package server

import (
	"net"
	"fmt"
	"log"
)


type IMServer struct {
	listener    net.Listener        // 服务端监听器 监听xx端口
	joinsniffer chan net.Conn       // 访问连接嗅探器 触发创建客户端连接处理方法
	//quitsniffer chan *common.Client // 连接退出嗅探器 触发连接退出处理方法

}


func StartIMServer()  {
	server := &IMServer{
		joinsniffer:make(chan net.Conn),
	}
	server.start()
}

func (server *IMServer)listen()  {
	go func() {
		for{
			select {
			case conn := <- server.joinsniffer
				server.joinHandler(conn)
			}
		}
	}()
}

func (server *IMServer)joinHandler(conn net.Conn)  {



}





func (server *IMServer)start()  {
	addr := fmt.Sprintf("0.0.0.0:%d", "1333")
	server.listener ,_ = net.Listen("tcp",addr)
	defer server.listener.Close()

	for{
		conn,err := server.listener.Accept()
		if err != nil {
			continue
		}
		log.Printf("新连接地址为:[%s]", conn.RemoteAddr())
		server.joinsniffer <- conn
	}
}


