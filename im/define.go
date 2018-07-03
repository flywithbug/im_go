package im

const (

	// heartbeat
	OP_HEARTBEAT     = int32(2)
	OP_HEARTBEAT_ACK = int32(3)


	//其他消息类型，都在客户端自行实现。
	// send  messgae
	OP_MSG     = int32(4)
	OP_MSG_ACK = int32(5) //客户端返回的ack 用于更新发送状态

	OP_MSG_SYNC = int32(6)    //消息发送者，发送消息给其他登录端

	OP_MSG_OFFLINE     = int32(7)  //离线消息

	OP_MSG_ROOM     = int32(10)  //离线消息


	// auth user
	OP_AUTH     = int32(18)
	OP_AUTH_ACK = int32(19)


)

func OperationMsg(operation int32) string {
	switch operation {
	case OP_AUTH:
		return "初次连接授权"
	}
	return "default_wait set notify info"
}
