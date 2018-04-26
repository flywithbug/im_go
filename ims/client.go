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



/*
 客户端列表
 */
type ClientTable map[string]*Client

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

func (client *Client)read() {
	for{
		if err := client.proto.ReadTCP(client.reader);err != nil {
			if err == io.EOF {
				client.close()
				fmt.Println("error",err)
				break
			}
			fmt.Println("else",err)
		}else {
			fmt.Println(client.proto,string(client.proto.Body),client.proto.Operation)
		}

	}
}

func (client *Client)write() {



}

//client 长链失败，关闭连接，处理数据保存事宜
func (client *Client)close()  {
	client.Conn.Close()
}




func (client *Client)Listen()  {
	go client.read()
	go client.write()
}


