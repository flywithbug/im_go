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
	listenAddr := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("初始化失败", err.Error())
		return
	}
	for {
		client, err := listen.Accept()
		if err != nil {
			return
		}
		go handleConnection(client)
	}
}

func handleConnection(conn net.Conn)  {
	client := NewClient(conn)
	client.Listen()
}