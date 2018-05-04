package im

import (
	log "github.com/flywithbug/log4go"
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
	if addrTcp, ok := addr.(*net.TCPAddr); ok {
		ip4 := addrTcp.IP.To4()
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
	case OP_SEND_MSG_ACK: //客户端返回的ack
		client.HandleACK(pro)
	case OP_SEND_MSG_SYNC_ACK: //客户端返回的ack
		client.HandleSyncACK(pro)
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
			log.Error(err.Error())
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
				log.Debug("client: %s socket closed", client.userId)
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

	atomic.AddInt64(&serverSummary.nconnections, -1)
	if client.uid > 0 {
		atomic.AddInt64(&serverSummary.nclients, -1)
	}

	atomic.StoreInt32(&client.closed, 1)
	log.Debug("close client:%d, %s", client.uid, client.userId)

	//client.RoomClient.Logout()
	//client.IMClient.Logout()
	client.RemoveClient()
	client.wt <- nil
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}
