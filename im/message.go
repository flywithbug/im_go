package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MessageIMInterface interface {
	ToData() []byte
	FromData(buff []byte) bool
}

const RawMessageHeaderLen = 20

//此模型只作为消息转发，存储解析 在model库中操作
type Message struct {
	sender 		int32   //4
	receiver 	int32
	timestamp 	int64   //8
	msgId		int32
	body 		[]byte
}
func (msg *Message) Description() string {
	return fmt.Sprintf("sender:%d,receiver:%d,timestamp:%d,msgId:%d,body:%s",msg.sender,msg.receiver,
		msg.timestamp,msg.msgId,msg.body)
}

func (msg *Message) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, msg.sender)
	binary.Write(buffer, binary.BigEndian, msg.receiver)
	binary.Write(buffer, binary.BigEndian, msg.timestamp)
	binary.Write(buffer, binary.BigEndian, msg.msgId)
	buffer.Write(msg.body)
	buf := buffer.Bytes()
	return buf
}

func (msg *Message) FromData(buff []byte) bool {
	if len(buff) < RawMessageHeaderLen {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &msg.sender)
	binary.Read(buffer, binary.BigEndian, &msg.receiver)
	binary.Read(buffer, binary.BigEndian, &msg.timestamp)
	binary.Read(buffer, binary.BigEndian, &msg.msgId)
	msg.body = buff[RawMessageHeaderLen:]
	return true
}










