package im

import (
	"fmt"
	log "github.com/flywithbug/log4go"
	"net"
)

//client 缓存
var appRoute *AppRoute

var serverSummary *ServerSummary

func init() {
	appRoute = NewAppRoute()
	serverSummary = NewServerSummary()
}

func listen(port int) {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	addr, _ := net.ResolveTCPAddr("tcp", address)
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("初始化失败", err.Error())
		return
	}
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			return
		}
		log.Info("新连接地址为:[%s] \n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	client := NewClient(conn)
	client.Listen()
}
