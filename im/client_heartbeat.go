package im

import (
	log "github.com/flywithbug/log4go"
)

//处理心跳包 回执
func (client *Client)HandleHeartbeat(pro *Proto)  {
	log.Debug("heartBeet")
	pro.Operation = OP_HEARTBEAT_ACK
	client.EnqueueMessage(*pro)
}