package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
	"fmt"
)

type ClientIM struct {
	*Connection
}

func (client *ClientIM)handleMessage(pro *Proto)  {
	switch pro.Operation {
	case OP_SEND_MSG:
		client.HandleIMMessage(pro)
	}

}

func (client *ClientIM)HandleIMMessage(pro *Proto)  {
	if client.uid == 0 {
		log.Warn("client has't been authenticated")
		return
	}
	var msg Message
	if !msg.FromData(pro.Body) {
		log.Warn("message decode not right")
		return
	}
	if msg.sender != client.uid {
		log.Warn("im message sender:%d client uid:%d\n", msg.sender, client.uid)
		return
	}
	msgId, err := model.SaveIMMessage(msg.sender,msg.receiver,0,msg.body)
	if err != nil {
		log.Warn(err.Error()+"消息存储服务出错")
		return
	}
	msg.msgId = msgId


	
}

func (client *ClientIM)handleImMessageACK(msg *Message,ver int16)  {

}