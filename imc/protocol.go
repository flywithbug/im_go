package imc

import (
	"errors"
	"encoding/json"
	"bytes"
	"encoding/binary"
	"io"
	log "github.com/flywithbug/log4go"
	"fmt"
)

type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
	SeqId     int32           `json:"seq"`  // sequence number chosen by client
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte) //解析对象
}

func (pro *Proto)Description()string{
	return fmt.Sprintf("ver:%d, operation:%d,seqId:%d,body:%s",pro.Ver,pro.Operation,pro.SeqId,pro.Body)
}

//
//const (
//	PackSize  	= 4  	//消息最大长度9kb
//	HeaderSize 	= 2
//	VerSize  	= 2
//	OperationSize = 4
//	SeqIdSize     = 4
//	RawHeaderSize = PackSize + HeaderSize + VerSize + OperationSize + SeqIdSize
//	MaxPackSize   = MaxBodySize + int32(RawHeaderSize)
//
//	// offset
//	PackOffset      = 0
//	HeaderOffset    = PackOffset + PackSize
//	VerOffset       = HeaderOffset + HeaderSize
//	OperationOffset = VerOffset + VerSize
//	SeqIdOffset     = OperationOffset + OperationSize
//)

//消息头 结构体
type protoHeader struct {
	bodyLen   int32 // 4 消息长度
	headerLen int16 // 2  //默认 RawHeaderSize = 16
	ver       int16 // 2
	op        int32 // 4
	seq       int32 // 4
}





// for tcp
const (
	MaxBodySize = int32(9 << 10)  //数据最大长度9kb
	RawHeaderSize = int16(16) //4+2+2+4+4
)

var(
	emptyProto = Proto{}
	emptyJSONBody = []byte("{}")

	ErrProtoPackLen   = errors.New("default server codec pack length error")
	ErrProtoHeaderLen = errors.New("default server codec header length error")

)

func WriteHeader(ph protoHeader,buffer io.Writer)  {
	//bodyLen 	int32   //消息长度
	//headerLen 	int16
	//ver 		int16
	//op			int32
	//seq			int32
	if ph.headerLen == 0 {
		ph.headerLen = RawHeaderSize
	}
	binary.Write(buffer,binary.BigEndian,ph.bodyLen)
	binary.Write(buffer,binary.BigEndian,ph.headerLen)
	binary.Write(buffer,binary.BigEndian,ph.ver)
	binary.Write(buffer,binary.BigEndian,ph.op)
	binary.Write(buffer,binary.BigEndian,ph.seq)
}


func ReadHeader(buff []byte)(*protoHeader,error)  {
	var ph protoHeader
	//var (
	//	bodyLen 	int32   //消息长度
	//	headerLen 	int16
	//	ver 		int16
	//	op			int32
	//	seq			int32
	//)
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer,binary.BigEndian,&ph.bodyLen)
	binary.Read(buffer,binary.BigEndian,&ph.headerLen)
	binary.Read(buffer,binary.BigEndian,&ph.ver)
	binary.Read(buffer,binary.BigEndian,&ph.op)
	binary.Read(buffer,binary.BigEndian,&ph.seq)
	if  ph.headerLen != RawHeaderSize{
		return nil,ErrProtoHeaderLen
	}
	return &ph,nil
}



func WriteMessage(w *bytes.Buffer, p *Proto)  {
	var ph protoHeader
	ph.bodyLen = int32(len(p.Body))
	ph.ver = p.Ver
	ph.op = p.Operation
	ph.seq = p.SeqId
	ph.headerLen = RawHeaderSize
	WriteHeader(ph,w)
	w.Write(p.Body)
}


func SendMessage(conn io.Writer,pro *Proto)error  {
	buffer := new(bytes.Buffer)
	WriteMessage(buffer,pro)
	buf := buffer.Bytes()
	n , err := conn.Write(buf)
	if err != nil {
		log.Info("sock write error:", err)
		return err
	}
	if n != len(buf) {
		log.Info("write less:%d %d", n, len(buf))
		return errors.New("write less")
	}
	return nil
}


func ReceiveLimitMessage(conn io.Reader,limitSize int)(pro *Proto)  {
	buff := make([]byte,RawHeaderSize)
	_,err := io.ReadFull(conn,buff)
	log.Info("receive header",buff)
	if err != nil {
		log.Info("sock read error:", err)
		return nil
	}
	ph,err := ReadHeader(buff)
	if err != nil {
		log.Info("buff read error:", err)
		return nil
	}
	if ph.bodyLen < 0 || int(ph.bodyLen) > limitSize {
		log.Info("invalid len:",ph.bodyLen)
		return nil
	}
	buff = make([]byte,ph.bodyLen)
	_,err = io.ReadFull(conn,buff)
	if err != nil {
		log.Info("sock read error:", err)
		return nil
	}
	pro = &emptyProto

	pro.Ver = ph.ver
	pro.SeqId = ph.seq
	pro.Operation = ph.op
	pro.Body = buff
	return pro
}


func ReceiveMessage(conn io.Reader) *Proto {
	return ReceiveLimitMessage(conn, 32*1024)
}

//消息大小限制在1M
func ReceiveStorageMessage(conn io.Reader) *Proto {
	return ReceiveLimitMessage(conn, 1024*1024)
}




const (
	// heartbeat
	OP_HEARTBEAT     = int32(2)
	OP_HEARTBEAT_ACK = int32(3)


	//其他消息类型，都在客户端自行实现。
	// send  messgae
	OP_MSG     = int32(4)
	OP_MSG_ACK = int32(5) //客户端返回的ack 用于更新发送状态

	OP_MSG_SYNC = int32(6)    //消息发送者，发送消息给其他登录端


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
