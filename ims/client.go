package ims

import (
	"net"
	"im_go/libs/bufio"
	"im_go/libs/proto"
	"fmt"
)

type Client struct {
	// 连接信息
	Key    	string        	//客户端连接的唯标志
	Conn   	net.TCPConn     //连接
	reader 	*bufio.Reader	//读取
	writer 	*bufio.Writer	//输出
}

func CreateClient(key string,conn net.TCPConn)*Client  {
	reader := bufio.NewReader(&conn)
	writer := bufio.NewWriter(&conn)
	client := &Client{
		Key:key,
		Conn:conn,
		reader:reader,
		writer:writer,
	}
	return client
}

func (client *Client)read()  {
	 for{
	 	var p proto.Proto
	 	err := p.ReadTCP(client.reader)
	 	fmt.Println(p,err)
	 }
}


