package im


//思路：1.收到群消息，把client add 到对应的roomSet 中，client退出时，移除
//思路：2.给roomId的消息绑定一个通知，收到这个消息，发送给监听的client去处理消息,roomId绑定了一群client客户端
type ClientROOM struct {
	*Connection
	roomId    int64
}