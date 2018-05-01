package imc

import (
	"bufio"
	"log"
	"net"
	"os"
	"fmt"
	"im_go/im"
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
			switch msg.Operation {
			case OP_AUTH_REPLY:
				var auth im.AuthenticationStatus
				auth.FromData(msg.Body)
				fmt.Println("授权状态",auth.Status)
			}
			fmt.Println("receive Msg",string(msg.Body))
		}
	}()

	for {
		line, _, _ := in.ReadLine()
		p := new(Proto)

		if string(line) == "auth"{
			var auth im.AuthenticationToken
			auth.Token = "6bde541f-1eb9-4600-a47e-8d7db7c7b460"
			auth.DeviceId = "4c6aba79-f768-4e26-8344-aa2b7bc173ec"
			auth.PlatformType = 3
			p.Ver = 1
			p.Body = auth.ToData()
			//if err !=nil {
			//	fmt.Println("Marshal",err)
			//}
			p.Operation = OP_AUTH
			p.SeqId = 1
		}else {
			p.Ver = 1
			p.Body = line
			p.Operation = OP_SEND_MSG
			p.SeqId = 1
		}

		err = SendMessage(conn,p)
		if err != nil {
			break
		}
		fmt.Println("send Msg",string(p.Body),err)
	}
}
