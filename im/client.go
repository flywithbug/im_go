package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type Client struct {
	Connection
	*ClientIM
	publicIp int32
}

func NewClient(conn *net.TCPConn) *Client {
	client := new(Client)
	client.conn = conn
	addr := conn.LocalAddr()
	if taddr, ok := addr.(*net.TCPAddr); ok {
		ip4 := taddr.IP.To4()
		if len(ip4) >= 4 {
			client.publicIp = int32(ip4[0])<<24 | int32(ip4[1])<<16 | int32(ip4[2])<<8 | int32(ip4[3])
		}
	}
	atomic.AddInt64(&serverSummary.nconnections, 1)
	client.wt = make(chan *Proto, 100)

	//消息处理器
	client.ClientIM = &ClientIM{&client.Connection}

	return client
}

func (client *Client) handleMessage(pro *Proto) {
	switch pro.Operation {
	case OP_AUTH:
		client.HandleAuthToken(pro)
	case OP_HEARTBEAT:
		//TODO
	}
	client.ClientIM.handleMessage(pro)

}

func (client *Client) AddClient() {
	route := appRoute.FindOrAddRoute(client.appid)
	route.AddClient(client)
}

func (client *Client) RemoveClient() {
	route := appRoute.FindRoute(client.appid)
	if route == nil {
		log.Warn("can't find app route")
		return
	}
	route.RemoveClient(client)

	//if client.room_id > 0 {
	//	route.RemoveRoomClient(client.room_id, client)
	//}
}

func (client *Client) Read() {
	for {
		tc := atomic.LoadInt32(&client.tc)
		if tc > 0 {
			log.Info("quit read goroutine, client:%d write goroutine blocked", client.uid)
			client.handleClientClosed()
			break
		}
		t1 := time.Now().Unix()
		msg, err := client.read()
		if err == io.EOF {
			client.handleClientClosed()
			break
		}
		if msg == nil {
			client.handleClientClosed()
			break
		}
		t2 := time.Now().Unix()
		if t2-t1 > 6*60-1 {
			log.Info("client:%d socket read timeout:%d %d", client.uid, t1, t2)
		}

		client.handleMessage(msg)
		t3 := time.Now().Unix()
		if t3-t2 > 2 {
			log.Info("client:%d handle message is too slow:%d %d", client.uid, t2, t3)
		}
	}

}

func (client *Client) Write() {
	running := true
	seq := 0

	for running {
		select {
		case pro := <-client.wt:
			if pro == nil {
				client.close()
				running = false
				log.Info("client: %s socket closed", client.userId)
				break
			}
			if pro.Operation == OP_SEND_MSG {
				atomic.AddInt64(&serverSummary.out_message_count, 1)
			}
			seq++
			pro.SeqId = int32(seq)

			client.send(pro)
		}
	}
	//等待200ms,避免发送者阻塞
	t := time.After(200 * time.Millisecond)
	running = true
	for running {
		select {
		case <-t:
			running = false
		case <-client.wt:
			log.Warn("msg is dropped")
			//case <- client.ewt:
			//	log.Warning("emsg is dropped")
		}
	}
}

func (client *Client) handleClientClosed() {
	close := atomic.LoadInt32(&client.closed)
	if client.uid > 0 {
		atomic.AddInt64(&serverSummary.nclients, -1)
		log.Info("HandleClientClosed client:%d %s", client.uid, client.userId)
	}
	if close == 0 && client.uid != 0 {
		atomic.AddInt64(&serverSummary.nconnections, -1)
		//quit when write goroutine received
		log.Info("close client userId:%s uid:%d", client.userId, client.uid)
		//client.RoomClient.Logout()
		//client.IMClient.Logout()
		client.RemoveClient()
		//数据库登录态置为0 ）
		model.UpdateUserStatus(client.uid, 0)

		client.uid = 0
	}
	client.wt <- nil
	atomic.StoreInt32(&client.closed, 1)
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}
