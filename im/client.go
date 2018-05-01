package im

import (
	"net"
	"sync/atomic"
	"fmt"
)

var clients map[string]*Client
type Client struct {
	Connection
	publicIp int32

}

func NewClient(conn net.Conn)*Client  {
	client := new(Client)
	client.conn = conn

	addr := conn.LocalAddr()
	if taddr, ok := addr.(*net.TCPAddr); ok {
		ip4 := taddr.IP.To4()
		if len(ip4) >= 4 {
			client.publicIp = int32(ip4[0]) << 24 | int32(ip4[1]) << 16 | int32(ip4[2]) << 8 | int32(ip4[3])
		}
	}
	atomic.AddInt64(&serverSummary.nconnections, 1)
	clients["a"] = client
	return client
}


func (client *Client)Read()  {
	fmt.Println("Lisetn Read")
	for {
		//tc := atomic.LoadInt32(&client.tc)
		//if tc > 0 {
		//	log.Infof("quit read goroutine, client:%d write goroutine blocked", client.uid)
		//	client.HandleClientClosed()
		//	break
		//}
		//t1 := time.Now().Unix()
		fmt.Println("read")
		msg := client.read()
		fmt.Println("receiveMSg",msg)
		//t2 := time.Now().Unix()
		//if t2 - t1 > 6*60 {
		//	log.Infof("client:%d socket read timeout:%d %d", client.uid, t1, t2)
		//}
		if msg == nil {
			client.HandleClientClosed()
			break
		}
		client.handleMessage(msg)
		//t3 := time.Now().Unix()
		//if t3 - t2 > 2 {
		//	log.Infof("client:%d handle message is too slow:%d %d", client.uid, t2, t3)
		//}

	}
}

func (client *Client)handleMessage(pro *Proto)  {
	fmt.Println("msg Operation",pro.Operation,OperationMsg(pro.Operation))
}







func (client *Client)Write()  {

}














func (client *Client) HandleClientClosed() {
	atomic.AddInt64(&serverSummary.nconnections, -1)


	//if client.uid > 0 {
	//	atomic.AddInt64(&server_summary.nclients, -1)
	////}
	//atomic.StoreInt32(&client.closed, 1)
	//
	//client.RemoveClient()
	//
	////quit when write goroutine received
	//client.wt <- nil
	//
	//client.RoomClient.Logout()
	//client.IMClient.Logout()
}



func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}
