package ims

import (
	"fmt"
	"net"
	"im_go/libs/bufio"
	"im_go/libs/proto"
)

var clientTables map[string]Client

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
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn)  {
	rr := bufio.NewReaderSize(conn,int(proto.MaxBodySize))
	p := new(proto.Proto)
	err := p.ReadTCP(rr)
	fmt.Println(p,string(p.Body),"\n",err)

}








