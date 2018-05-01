package im

import (
	"errors"
	"encoding/json"
	"bytes"
	 "encoding/binary"
	"io"
)

type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
	SeqId     int32           `json:"seq"`  // sequence number chosen by client
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)
}



//消息头 结构体
type protoHeader struct {
	packLen 	int32   // 4 消息长度
	headerLen	int16	// 2
	ver 		int16	// 2
	op			int32	// 4
	seq			int32	// 4

	//并不插入结构体
	bodyLen 	int    //等于 packLen-headerLen
}




// for tcp
const (
	MaxBodySize = int32(9 << 10)  //数据最大长度9kb
	RawHeaderSize = int16(16) //4+2+2+4+4
)
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

var(
	emptyProto = Proto{}
	emptyJSONBody = []byte("{}")

	ErrProtoPackLen   = errors.New("default server codec pack length error")
	ErrProtoHeaderLen = errors.New("default server codec header length error")

	ProtoReady 	= &Proto{Operation:OP_PROTO_READY}
	ProtoFinish = &Proto{Operation:OP_PROTO_FINISH}
)

func WriteHeader(ph protoHeader,buffer io.Writer)  {
	//packLen 	int32   //消息长度
	//headerLen 	int16
	//ver 		int16
	//op			int32
	//seq			int32
	if ph.headerLen == 0 {
		ph.headerLen = RawHeaderSize
	}
	binary.Write(buffer,binary.BigEndian,ph.packLen)
	binary.Write(buffer,binary.BigEndian,ph.headerLen)
	binary.Write(buffer,binary.BigEndian,ph.ver)
	binary.Write(buffer,binary.BigEndian,ph.op)
	binary.Write(buffer,binary.BigEndian,ph.seq)
}


func ReadHeader(buff []byte)(*protoHeader,error)  {
	var ph protoHeader
	//var (
	//	packLen 	int32   //消息长度
	//	headerLen 	int16
	//	ver 		int16
	//	op			int32
	//	seq			int32
	//)

	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer,binary.BigEndian,&ph.packLen)
	binary.Read(buffer,binary.BigEndian,&ph.headerLen)
	binary.Read(buffer,binary.BigEndian,&ph.ver)
	binary.Read(buffer,binary.BigEndian,&ph.op)
	binary.Read(buffer,binary.BigEndian,&ph.seq)
	if  ph.headerLen != RawHeaderSize{
		return nil,ErrProtoHeaderLen
	}
	return &ph,nil

}




