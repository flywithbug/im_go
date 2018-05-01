package im

import (
	"net"
	"sync/atomic"
	log "github.com/flywithbug/log4go"
	"time"
)

type Client struct {
	Connection
	publicIp int32

}

func NewClient(conn *net.TCPConn)*Client  {
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
	client.out = make(chan *Proto,100)
	return client
}

func (client *Client)Read()  {
	for {
		tc := atomic.LoadInt32(&client.tc)
		if tc > 0 {
			log.Info("quit read goroutine, client:%d write goroutine blocked", client.uid)
			client.HandleClientClosed()
			break
		}
		t1 := time.Now().Unix()
		msg := client.read()
		t2 := time.Now().Unix()
		if t2 - t1 > 6*60 {
			log.Info("client:%d socket read timeout:%d %d", client.uid, t1, t2)
		}
		if msg == nil {
			client.HandleClientClosed()
			break
		}
		client.handleMessage(msg)
		t3 := time.Now().Unix()
		if t3 - t2 > 2 {
			log.Info("client:%d handle message is too slow:%d %d", client.uid, t2, t3)
		}

	}
}

func (client *Client)handleMessage(pro *Proto)  {
	//fmt.Println("receiveMSg:",string(pro.Body),len(pro.Body),pro.Operation)
	switch pro.Operation {
	case OP_AUTH:
		client.HandleAuthToken(pro)
	case OP_SEND_MSG_REPLY:
		//消息回执
	}

	//client.out <- pro

}


func (client *Client)Write()  {
	running := true
	seq := 0

	for running {
		select {
		case pro := <- client.out:
			if pro == nil{
				client.close()
				running = false
				log.Info("client:%d socket closed", client.uid)
				break
			}
			if pro.Operation == OP_SEND_MSG {
				atomic.AddInt64(&serverSummary.out_message_count,1)
			}
			seq++
			//p := &Proto{Ver:pro.Ver,}
			client.send(pro)
		}
	}
	//等待200ms,避免发送者阻塞
	//t := time.After(200 * time.Millisecond)
	//running = true
	//for running {
	//	select {
	//	case <- t:
	//		running = false
	//	//case <- client.wt:
	//	//	log.Warning("msg is dropped")
	//	//case <- client.ewt:
	//	//	log.Warning("emsg is dropped")
	//	}
	//}

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
