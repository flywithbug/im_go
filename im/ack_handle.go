package im

import (
	log "github.com/flywithbug/log4go"
)

func (client *Client) HandleACK(pro *Proto) {
	//TODO 完善回执机制
	log.Info("ack:", pro.SeqId)
}
