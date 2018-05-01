package imc

import (
	"bufio"
	"log"
	"net"
	"os"
	"fmt"
)

func StartClient(port int) {
	address := fmt.Sprintf("0.0.0.0:%d",port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	//reader := bufio.NewReader(conn)
	in := bufio.NewReader(os.Stdin)
	go func() {
		for {
			msg := ReceiveMessage(conn)
			if msg == nil{
				break
			}
			fmt.Println("receive Msg",msg)
		}
	}()

	for {
		line, _, _ := in.ReadLine()
		// 模拟一个请求
		// {"command":"GET_CONN","data":null}
		// {"command":"GET_BUDDY_LIST","data":null}
		//buffer := new(bytes.Buffer)
		//var ph protoHeader
		//ph.headerLen = RawHeaderSize
		//ph.seq = 1
		//ph.op = 2
		//ph.bodyLen = int32(len(line))
		//ph.ver = 1
		//WriteHeader(ph,buffer)
		//buffer.Write(line)
		//bb := buffer.Bytes()
		//writer.Write(bb)
		p := new(Proto)
		p.Ver = 1
		p.Body = line
		p.Operation = OP_AUTH
		p.SeqId = 1
		err = SendMessage(conn,p)
		fmt.Println("send Msg",p,err)
	}
}
