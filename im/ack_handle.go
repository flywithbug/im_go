package im

import (
	log "github.com/flywithbug/log4go"
	"im_go/model"
)



/*
// 客户端返回的ack 用于标记消息已发送成功
*/
func (client *Client) handleMessageACK(pro *Proto) {
	//TODO 优化为rpc和方式修改  通知发发送方，送达，（send operation）OP_MSG_DELIVER_ACK
	var ack MessageACK
	ack.FromData(pro.Body)
	err := model.SaveIMMessageFromOfflineMessage(ack.msgId)
	//err := model.UpdateMessageACK(ack.msgId,1)
	if err != nil {
		log.Error("error"+err.Error() + ack.Description())
	}

}

//已读回执
func (client *Client)handleMessageReadAck(pro *Proto)  {
	var ack MessageACK
	ack.FromData(pro.Body)
	err := model.UpdateMessageACK(ack.msgId,2)
	if err != nil {
		log.Error("error"+err.Error() + ack.Description())
	}
}




func (client *Client) handleSyncACK(pro *Proto) {
	//多客户端同步消息回执
	//var ack MessageACK
	//ack.FromData(pro.Body)
	////err := model.UpdateMessageACK(ack.msgId,1)
	//if err != nil {
	//}
}





