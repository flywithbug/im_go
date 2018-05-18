package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
)

/*
// 客户端返回的ack 用于标记消息已发送成功 //send delivryAck to sender
*/
func (client *Client) handleMessageACK(pro *Proto) {
	//TODO 优化为rpc和方式修改
	var ack MessageACK
	ack.FromData(pro.Body)
	err := model.UpdateMessageACK(ack.msgId,1)
	if err != nil {
		log.Error("error"+err.Error() + ack.Description())
	}
}

//已读回执
func (client *Client)handleMessageReadAck(pro *Proto)  {

}



func (client *Client) handleSyncACK(pro *Proto) {
	//多客户端同步消息回执
	//var ack MessageACK
	//ack.FromData(pro.Body)
	////err := model.UpdateMessageACK(ack.msgId,1)
	//if err != nil {
	//}
}





