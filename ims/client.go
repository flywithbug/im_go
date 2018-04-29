package ims

import (
	"net"
	"im_go/libs/bufio"
	"im_go/libs/proto"
	"fmt"
	"io"
	"im_go/libs/define"
	"sync/atomic"
)

type Client struct {
	Connection
	// 连接信息
	Key    	string        	//客户端连接的唯标志
	//reader 	*bufio.Reader	//读取
	//writer 	*bufio.Writer	//输出
	pro   	*proto.Proto
}


func NewClient(conn *net.TCPConn)(client *Client)  {
	client = new(Client)
	client.bufPrepare()
	reader := bufio.NewReaderSize(conn,int(proto.MaxPackSize))
	writer := bufio.NewWriter(conn)
	p := new(proto.Proto)
	client.conn = conn
	client.pro = p
	client.reader = reader
	client.writer = writer
	atomic.AddInt64(&serverSummary.nConnections, 1)

	return client
}

func (client *Client)Read() {
	defer client.conn.Close()
	for{
		if msg,err := client.read();err != nil{
			if err == io.EOF {
				fmt.Println("close:",err)
			}
			fmt.Println("else:",err)
			break
		}else {
			go client.handleMessage(msg)
		}
	}
	client.HandleClientClosed()
}

func (client *Client)Write() {



}


//消息处理
func (client *Client)handleMessage(pro *Message)  {
	client.count++
	fmt.Println("count",client.count,"操作类型",pro.Operation,string(pro.Body))
	switch pro.Operation {
	case define.OP_AUTH:
		//fmt.Println(pro,string(pro.Body),pro.Operation)
	case define.OP_PROTO_FINISH:
		client.HandleClientClosed()

	}
}

//client 长链失败，关闭连接，处理数据保存事宜
func (client *Client)HandleClientClosed()  {
	fmt.Println("client close",client.Key)
	atomic.AddInt64(&serverSummary.nConnections, -1)
	client.close()

}

func (client *Client) AddClient() {
	route := appRoute.FindOrAddRoute(client.appId)
	route.AddClient(client)
}

func (client *Client) RemoveClient() {
	route := appRoute.FindRoute(client.appId)
	if route == nil {
		//log.Warning("can't find app route")
		return
	}
	route.RemoveClient(client)
}





func (client *Client)Listen()  {
	go client.Read()
	go client.Write()
}


