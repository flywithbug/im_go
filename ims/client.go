package ims

import (
	"net"
	"fmt"
	"io"
	"im_go/libs/define"
	"sync/atomic"
)

type Client struct {
	Connection
	// 连接信息
	Key    	string        	//客户端连接的唯标志
	publicIp int32
}


func NewClient(conn *net.TCPConn)(client *Client)  {
	client = new(Client)
	client.bufPrepare(conn)
	conn.LocalAddr()
	addr := conn.LocalAddr()
	if ad, ok := addr.(*net.TCPAddr); ok {
		ip4 := ad.IP.To4()
		client.publicIp = int32(ip4[0]) << 24 | int32(ip4[1]) << 16 | int32(ip4[2]) << 8 | int32(ip4[3])
	}
	atomic.AddInt64(&serverSummary.nConnections, 1)
	return client
}

func (client *Client)Read() {
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
func (client *Client)handleMessage(msg *Message)  {
	//time.Sleep(time.Microsecond*20)
	client.count++
	fmt.Println("count",client.count,"操作类型",msg.Operation,string(msg.Body))
	switch msg.Operation {
	case define.OP_AUTH:
		auth := new(AuthenticationToken)
		auth.FromData(msg.Body)
		client.HandleAuthToken(auth,msg.Ver)
	case define.OP_PROTO_FINISH:
		client.HandleClientClosed()

	}
}

func (client *Client)HandleAuthToken(login *AuthenticationToken,version int16)  {
	//fmt.Println("auth",login)

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


