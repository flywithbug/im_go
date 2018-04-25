package proto

import (
	"testing"
	"im_go/libs/bytes"
	"fmt"
	"im_go/libs/encoding/binary"


)

func TestProto_WriteTo(t *testing.T) {
	body := []byte("hello,world,hello,world,hello,world,hello,world")

	p := Proto{
		Ver:1111,
		Operation:2222,
		SeqId:333,
		Body:body,
	}

	var w  bytes.Writer
	p.WriteTo(&w)
	buf := w.Buffer()
	fmt.Println(buf,w.Size(),string(buf))
	var p1 Proto

	packLen := binary.BigEndian.Int32(buf[PackOffset:HeaderOffset])
	headerLen := binary.BigEndian.Int16(buf[HeaderOffset:VerOffset])
	p1.Ver = binary.BigEndian.Int16(buf[VerOffset:OperationOffset])
	p1.Operation = binary.BigEndian.Int32(buf[OperationOffset:SeqIdOffset])
	p1.SeqId = binary.BigEndian.Int32(buf[SeqIdOffset:])

	fmt.Println(packLen,headerLen,p1.Ver,p1.Operation,p1.SeqId,string(p1.Body))



}

func unPackage([]byte)  {



}

func testLongMsg(line []byte)  {

	p := Proto{
		Ver:10,
		Operation:20,
		SeqId:902020,
		Body:line,
	}

	var w  bytes.Writer
	p.WriteTo(&w)
	buf := w.Buffer()
	fmt.Println(buf,w.Size(),string(buf))
	var p1 Proto

	packLen := binary.BigEndian.Int32(buf[PackOffset:HeaderOffset])
	headerLen := binary.BigEndian.Int16(buf[HeaderOffset:VerOffset])
	p1.Ver = binary.BigEndian.Int16(buf[VerOffset:OperationOffset])
	p1.Operation = binary.BigEndian.Int32(buf[OperationOffset:SeqIdOffset])
	p1.SeqId = binary.BigEndian.Int32(buf[SeqIdOffset:])
	fmt.Println(packLen,headerLen,p1.Ver,p1.Operation,p1.SeqId)
}