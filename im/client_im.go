package im

import (
	log "github.com/flywithbug/log4go"
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






}