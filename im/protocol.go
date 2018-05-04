package im

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/flywithbug/log4go"
	"io"
)

type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
	SeqId     int32           `json:"seq"`  // sequence number chosen by client
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte) //解析对象
}

func (pro *Proto) Description() string {
	return fmt.Sprintf("ver:%d, operation:%d,seqId:%d,body:%s", pro.Ver, pro.Operation, pro.SeqId, pro.Body)
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
	MaxBodySize   = int32(9 << 10) //数据最大长度9kb
	RawHeaderSize = int16(16)      //4+2+2+4+4
)

var (
	emptyProto    = Proto{}
	emptyJSONBody = []byte("{}")

	ErrProtoPackLen   = errors.New("default server codec pack length error")
	ErrProtoHeaderLen = errors.New("default server codec header length error")

	ProtoReady  = &Proto{Operation: OP_PROTO_READY}
	ProtoFinish = &Proto{Operation: OP_PROTO_FINISH}
)

func WriteHeader(ph protoHeader, buffer io.Writer) {
	//bodyLen 	int32   //消息长度
	//headerLen 	int16
	//ver 		int16
	//op			int32
	//seq			int32
	if ph.headerLen == 0 {
		ph.headerLen = RawHeaderSize
	}
	binary.Write(buffer, binary.BigEndian, ph.bodyLen)
	binary.Write(buffer, binary.BigEndian, ph.headerLen)
	binary.Write(buffer, binary.BigEndian, ph.ver)
	binary.Write(buffer, binary.BigEndian, ph.op)
	binary.Write(buffer, binary.BigEndian, ph.seq)
}

func ReadHeader(buff []byte) (*protoHeader, error) {
	var ph protoHeader
	//var (
	//	bodyLen 	int32   //消息长度
	//	headerLen 	int16
	//	ver 		int16
	//	op			int32
	//	seq			int32
	//)
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &ph.bodyLen)
	binary.Read(buffer, binary.BigEndian, &ph.headerLen)
	binary.Read(buffer, binary.BigEndian, &ph.ver)
	binary.Read(buffer, binary.BigEndian, &ph.op)
	binary.Read(buffer, binary.BigEndian, &ph.seq)
	if ph.headerLen != RawHeaderSize {
		return nil, ErrProtoHeaderLen
	}
	return &ph, nil
}

func WriteMessage(w *bytes.Buffer, p *Proto) {
	var ph protoHeader
	ph.bodyLen = int32(len(p.Body))
	ph.ver = p.Ver
	ph.op = p.Operation
	ph.seq = p.SeqId
	ph.headerLen = RawHeaderSize
	WriteHeader(ph, w)
	w.Write(p.Body)
}

func SendMessage(conn io.Writer, pro *Proto) error {
	buffer := new(bytes.Buffer)
	WriteMessage(buffer, pro)
	buf := buffer.Bytes()
	n, err := conn.Write(buf)
	if err != nil {
		log.Debug("sock write error:", err)
		return err
	}
	if n != len(buf) {
		log.Debug("write less:%d %d", n, len(buf))
		return errors.New("write less")
	}
	return nil
}

func ReceiveLimitMessage(conn io.Reader, limitSize int) (pro *Proto, err error) {
	buff := make([]byte, RawHeaderSize)
	_, err = io.ReadFull(conn, buff)
	if err != nil {
		log.Debug("sock read error:%s", err.Error())
		return nil, err
	}
	ph, err := ReadHeader(buff)
	if err != nil {
		log.Info("buff read error:%s", err.Error())
		return nil, err
	}
	if ph.bodyLen < 0 || int(ph.bodyLen) > limitSize {
		log.Info("invalid len:%d", ph.bodyLen)
		return nil, err
	}
	buff = make([]byte, ph.bodyLen)
	_, err = io.ReadFull(conn, buff)
	if err != nil {
		//log.Info("sock read error:%s", err.Error())
		return nil, err
	}
	pro = &emptyProto
	pro.Ver = ph.ver
	pro.SeqId = ph.seq
	pro.Operation = ph.op
	pro.Body = buff
	return pro, nil
}

func ReceiveMessage(conn io.Reader) (pro *Proto, err error) {
	return ReceiveLimitMessage(conn, 32*1024)
}

//消息大小限制在1M
func ReceiveStorageMessage(conn io.Reader) (pro *Proto, err error) {
	return ReceiveLimitMessage(conn, 1024*1024)
}
