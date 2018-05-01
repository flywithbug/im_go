package im

import (
	"fmt"
	"net"
)

func init()  {
	serverSummary = NewServerSummary()
	clients = make(map[string]*Client,10)
}

func Listen(port int)  {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	addr,_ := net.ResolveTCPAddr("tcp",address)
	listen, err := net.ListenTCP("tcp",addr)
	if err != nil {
		fmt.Println("初始化失败", err.Error())
		return
	}
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			return
		}
		fmt.Printf("新连接地址为:[%s]", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn)  {
	client := NewClient(conn)
	client.Listen()
}