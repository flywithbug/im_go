package im

import (
	log "github.com/flywithbug/log4go"
)



func (client *Client)HandleACK(pro *Proto)  {
	log.Info("ack:", pro.SeqId)
}