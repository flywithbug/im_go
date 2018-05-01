package im

import (
	"fmt"
	"net"
	"github.com/golang/glog"
)


func Listen(port int)  {
	address := fmt.Sprintf("0.0.0.0:%d",port)
	addr,_:=net.ResolveTCPAddr("tcp",address)
	listen ,err := net.ListenTCP("tcp",addr)
	if err != nil {
		glog.Info("初始化失败", err.Error())
		return
	}
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			return
		}
		glog.Infof("新连接地址为:[%s]", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn)  {
	client := NewClient(conn)
	client.Listen()
	go func() {
		client.SendMsg("connect Success")
	}()

}