package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
	"fmt"
	"sync/atomic"
)

type ClientIM struct {
	*Connection
}

func (client *ClientIM) handleMessage(pro *Proto) {
	switch pro.Operation {
	case OP_SEND_MSG:
		client.HandleIMMessage(pro)
	}

}

func (client *ClientIM) HandleIMMessage(pro *Proto) {
	if client.uid == 0 {
		log.Warn("client has't been authenticated")
		return
	}
	msg := new(Message)
	if !msg.FromData(pro.Body) {
		log.Error(fmt.Sprintf("message decode not right,body: %s,%d" , pro.Body ,client.uid))
		return
	}

	if msg.sender != client.uid {
		log.Warn("im message sender:%d client uid:%d\n", msg.sender, client.uid)
		return
	}

	//消息存入服务器
	msgId, err := model.SaveIMMessage(msg.sender, msg.receiver, 0, msg.body)
	if err != nil {
		log.Warn(err.Error() + "消息存储服务出错")
		return
	}

	msg.msgId = msgId
	pro.Body = msg.ToData()

	//发送消息给receiver
	client.SendMessage(msg.receiver, pro)
	//发送消息给其他登录登陆点
	pro.Operation = OP_SEND_MSG_SYNC
	client.SendMessage(msg.sender,pro)
	//消息回执
	client.handleImMessageACK(msgId, client.version, pro.SeqId)

	atomic.AddInt64(&serverSummary.in_message_count, 1)

}

func (client *ClientIM) handleImMessageACK(msgId int32, ver int16, seq int32) {
	ackMsg := new(MessageACK)
	ackMsg.msgId = msgId

	ack := Proto{}
	ack.Ver = ver
	ack.Operation = OP_SEND_MSG_ACK
	ack.Body = ackMsg.ToData()
	ack.SeqId = seq
	client.EnqueueMessage(ack)
	//客户端收到回执的msgId 才算消息发送完毕
}

func (client *ClientIM)sendOffLineMessage()  {
	//TODO 优化为rpc和方式获取
	ms,err :=  model.FindeMessagesReceiver(client.uid,0)
	if err != nil{
		log.Error(err.Error())
		return
	}
	p := Proto{}
	for _,imMsg := range ms{
		//fmt.Printf("offline msg :%s",imMsg.Description())
		p.Operation = OP_SEND_MSG
		p.Ver = client.version
		p.Body = FromIMMessage(&imMsg).ToData()
		p.SeqId = imMsg.Id
		client.EnqueueMessage(p)
	}
}

func FromIMMessage(imMsg *model.IMMessage)(msg *Message)  {
	msg = new(Message)
	msg.sender = imMsg.Sender
	msg.receiver = imMsg.Receiver
	msg.msgId = imMsg.Id
	msg.body = imMsg.Content
	msg.timestamp = imMsg.TimeStamp
	return msg
}