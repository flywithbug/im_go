package im

import (
	"net"
	"sync/atomic"
	log "github.com/golang/glog"
)

type Client struct {
	Connection

}

func NewClient(conn *net.TCPConn)*Client  {
	client := new(Client)
	client.conn = conn

	return client
}


func (client *Client)Read()  {
	for {
		tc := atomic.LoadInt32(&client.tc)
		if tc > 0 {
			log.Infof("quit read goroutine, client:%d write goroutine blocked", client.uid)
			client.HandleClientClosed()
			break
		}


	}
}


func (client *Client)Write()  {

}












func (client *Client) HandleClientClosed() {
	//atomic.AddInt64(&server_summary.nconnections, -1)
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
