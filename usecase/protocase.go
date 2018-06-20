package main


//import (
//	"fmt"
//	"im_go/libs/proto"
//	"im_go/libs/bytes"
//	"im_go/libs/encoding/binary"
//	"bufio"
//	"os"
//)
//
//func main()  {
//	in := bufio.NewReader(os.Stdin)
//	for  {
//		line, _, _ := in.ReadLine()
//		fmt.Println(line)
//		LongMsgInput(line)
//	}
//}
//
//func LongMsgInput(line []byte)  {
//	p := proto.Proto{
//		Ver:10,
//		Operation:20,
//		SeqId:902020,
//		Body:line,
//	}
//
//	var w  bytes.Writer
//	p.WriteTo(&w)
//	buf := w.Buffer()
//	fmt.Println(buf,w.Size(),string(buf))
//	var p1 proto.Proto
//
//	packLen := binary.BigEndian.Int32(buf[proto.PackOffset:proto.HeaderOffset])
//	headerLen := binary.BigEndian.Int16(buf[proto.HeaderOffset:proto.VerOffset])
//	p1.Ver = binary.BigEndian.Int16(buf[proto.VerOffset:proto.OperationOffset])
//	p1.Operation = binary.BigEndian.Int32(buf[proto.OperationOffset:proto.SeqIdOffset])
//	p1.SeqId = binary.BigEndian.Int32(buf[proto.SeqIdOffset:])
//	bodyLen := int(packLen - int32(headerLen))
//	p.Body = w.Peek(bodyLen)
//	fmt.Println("body",p.Body)
//	fmt.Println(packLen,headerLen,p1.Ver,p1.Operation,p1.SeqId,p.Body)
//}