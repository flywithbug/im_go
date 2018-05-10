package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
)

func (client *Client) HandleACK(pro *Proto) {

	//TODO 优化为rpc和方式修改
	var ack MessageACK
	ack.FromData(pro.Body)
	err := model.UpdateMessageACK(ack.msgId,1)
	if err != nil {
		log.Error("error"+err.Error() + ack.Description())
	}
}


func (client *Client) HandleSyncACK(pro *Proto) {
	//多客户端同步消息回执
	//var ack MessageACK
	//ack.FromData(pro.Body)
	////err := model.UpdateMessageACK(ack.msgId,1)
	//if err != nil {
	//}
}





