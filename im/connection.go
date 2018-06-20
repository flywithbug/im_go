package im

import (
	log "github.com/flywithbug/log4go"
	"net"
	"sync/atomic"

	"fmt"
	"time"
)

const CLIENT_TIMEOUT = (60 * 6)

type Connection struct {
	conn *net.TCPConn

	closed int32

	tc int32 //write channel timeout count
	wt chan *Proto

	tm time.Time

	//协议版本号
	version int16

	forbidden int32 //是否被禁言
	//notification_on bool //桌面在线时是否通知手机端
	online bool

	appId      int64  //登录用户所属AppId
	uid        int32  //登录用户的数据Id
	deviceId   string //设备唯一Id
	platformId int8   //设备类型Id
	userId     string //登录用户的UserId
	Token      string //登录用户的token
}

func (client *Connection) read() (*Proto, error) {
	client.conn.SetReadDeadline(time.Now().Add(CLIENT_TIMEOUT * time.Second))
	return ReceiveMessage(client.conn)
}

func (client *Connection) send(pro *Proto) {
	tc := atomic.LoadInt32(&client.tc)
	if tc > 0 {
		log.Info("can't write data to blocked socket")
		return
	}
	client.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
	err := SendMessage(client.conn, pro)
	if err != nil {
		atomic.AddInt32(&client.tc, 1)
		log.Info("send msg:", OperationMsg(pro.Operation), " tcp err:", err)
	}
}

// 根据连接类型关闭
func (client *Connection) close() {
	client.conn.Close()
}

//把消息加入到发送队列中 拷贝消息对象，不使用指针传递
func (client *Connection) EnqueueMessage(p Proto) bool {
	//warning 隔离指针传递
	//p := new(Proto)
	//p.SeqId = pro.SeqId
	//p.Body = pro.Body
	//p.Ver = pro.Ver
	//p.Operation = pro.Operation

	closed := atomic.LoadInt32(&client.closed)
	if closed > 0 {
		log.Info("can't send message to closed connection:%d %s", client.uid,client.userId)
		return false
	}
	tc := atomic.LoadInt32(&client.tc)
	if tc > 0 {
		log.Info("can't send message to blocked connection:%d", client.uid)
		atomic.AddInt32(&client.tc, 1)
		return false
	}

	select {
	case client.wt <- &p:
		return true
	case <-time.After(60 * time.Second):
		atomic.AddInt32(&client.tc, 1)
		log.Info("send message to wt timed out:%d", client.uid)
		return false
	}
}


func (client *Connection) SendMessage(uid int32, pro *Proto) bool {
	//发送推送给offline的uId
	appid := client.appId
	route := appRoute.FindRoute(appid)
	if route == nil {
		log.Error(fmt.Sprintf(" app not found: can't send message, appid:%d uid:%d cmd:%d", appid, uid, pro.Operation))
		return false
	}
	clients := route.FindClientSet(uid)
	if len(clients) == 0 {
		//走推送通道
		log.Debug(fmt.Sprintf("TODO(need offline push )can't send message, appid:%d uid:%d cmd:%d", appid, uid, pro.Operation))
		PushServiceHandler(uid,appid,pro)
		return false
	}

	//fmt.Printf("======clients===len:%d uid:%d \n",len(clients),uid)
	send := false
	for c := range clients {
		if &c.Connection == client {
			continue
		}
		if c.EnqueueMessage(*pro) {
			send = true
		}
	}
	return send
}

