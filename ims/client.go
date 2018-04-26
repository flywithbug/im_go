package ims

import (
	"net"
	"im_go/libs/bufio"
	"im_go/libs/proto"
	"fmt"
	"io"
)

type Client struct {
	// 连接信息
	Key    	string        	//客户端连接的唯标志
	Conn   	*net.TCPConn     //连接
	reader 	*bufio.Reader	//读取
	writer 	*bufio.Writer	//输出
	proto   *proto.Proto

							//补充连接信息
}

func NewClient(key string,conn *net.TCPConn)*Client  {
	reader := bufio.NewReaderSize(conn,int(proto.MaxPackSize))
	writer := bufio.NewWriter(conn)
	p := new(proto.Proto)
	client := &Client{
		Key:key,
		Conn:conn,
		reader:reader,
		writer:writer,
		proto:p,
	}
	fmt.Println("client init")
	return client
}

func (client *Client)read()  {
	for{
		if err := client.proto.ReadTCP(client.reader);err == nil {
			fmt.Println(client.proto,string(client.proto.Body),"\n",err)
		}else {
			if err == io.EOF {
				fmt.Println("error",err)
				break
			}
			//client.Conn.Close()
			fmt.Println("else",err)
		}
	 }
}

func (client *Client)Listen()  {
	go client.read()

	//rr := bufio.NewReaderSize(client.Conn,int(proto.MaxBodySize))
	//p := new(proto.Proto)
	//err := p.ReadTCP(rr)
	//fmt.Println(string(p.Body),"\n",err)

}


