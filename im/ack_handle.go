package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
)

func (client *Client) HandleACK(pro *Proto) {
	//TODO 客户端返回ack
	log.Info("ack:", pro.SeqId)
	var ack MessageACK
	ack.FromData(pro.Body)
	err := model.UpdateMessageACK(ack.msgId,1)
	if err != nil {
		log.Error("error"+err.Error())
	}
}
