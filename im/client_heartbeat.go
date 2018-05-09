package im


//处理心跳包 回执
func (client *Client)HandleHeartbeat(pro *Proto)  {
	pro.Operation = OP_HEARTBEAT_ACK

	client.EnqueueMessage(*pro)
}