package im


type MessageIMInterface interface {
	ToData() []byte
	FromData(buff []byte) bool
}


//此模型只作为消息转发，存储解析 在model库中操作
type Message struct {
	sender 		int64
	receiver 	int64
	content 	string
}







