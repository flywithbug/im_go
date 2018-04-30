package ims

import (
	"net"
	"fmt"
	"io"
	"im_go/libs/define"
	"sync/atomic"
	"im_go/model"
	"github.com/golang/glog"
	"time"
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
	addr := conn.LocalAddr()
	if ad, ok := addr.(*net.TCPAddr); ok {
		ip4 := ad.IP.To4()
		if len(ip4)>4 {
			client.publicIp = int32(ip4[0]) << 24 | int32(ip4[1]) << 16 | int32(ip4[2]) << 8 | int32(ip4[3])
		}
	}
	atomic.AddInt64(&serverSummary.nConnections, 1)
	client.wt = make(chan *Message,100)
	fmt.Println("new client ",client.publicIp)

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
	seq := 0
	running := true
	//loaded := false
	//for running && !loaded{
	//	//select {
	//	//
	//	//}
	//	loaded = true
	//}

	for running {
		select {
		case msg := <- client.wt:
			if msg == nil {
				client.close()
				running = false
				glog.Infof("client:%d socket closed", client.uId)
				break
			}
			if msg.Operation ==define.OP_SEND_MSG {
				atomic.AddInt64(&serverSummary.outMessageCount, 1)
			}
			seq++
			vMsg := new(Message)
			vMsg.Ver = msg.Ver
			vMsg.Operation = msg.Operation
			vMsg.SeqId = int32(seq)
			vMsg.Body = msg.Body
			err :=client.send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	//等待200ms,避免发送者阻塞
	t := time.After(200 *time.Millisecond)
	running = true
	for running {
		select {
		case <- t:
			running = false
		case <- client.wt:
			glog.Warning("msg is dropped")
		//case <- client.ewt:
		//	log.Warning("emsg is dropped")
		}
	}


}


//消息处理
func (client *Client)handleMessage(msg *Message)  {
	//time.Sleep(time.Microsecond*20)
	client.count++
	switch msg.Operation {
	case define.OP_AUTH:
		auth := new(AuthenticationToken)
		auth.FromData(msg.Body)
		client.HandleAuthToken(auth,msg.Ver)
	case define.OP_PROTO_FINISH:
		client.HandleClientClosed()
	default:
		fmt.Println("count",client.count,"操作类型",msg.Operation,string(msg.Body))
	}
}

func (client *Client)HandleAuthToken(auth *AuthenticationToken,version int16)  {
	login,err   := model.GetLoginByToken(auth.token,model.STATUS_LOGIN)
	if err != nil{
		fmt.Println("auth",err,login)
		msg := new(Message)
		msg.Reset()
		msg.Operation = define.OP_AUTH_REPLY
		msg.Ver = version
		authStatus := AuthenticationStatus{-1,0}//授权失败
		msg.Body = authStatus.ToData(version)
		client.EnqueueMessage(msg)
		return
	}else{
		fmt.Println(login)
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


