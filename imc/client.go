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
			fmt.Println("receive Msg",string(msg.Body))
		}
	}()

	for {
		line, _, _ := in.ReadLine()

		p := new(Proto)
		p.Ver = 1
		p.Body = line
		p.Operation = OP_AUTH
		p.SeqId = 1
		err = SendMessage(conn,p)
		if err != nil {
			break
		}
		fmt.Println("send Msg",string(p.Body),err)
	}
}
