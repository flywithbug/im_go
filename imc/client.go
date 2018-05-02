package imc

import (
	"bufio"
	"net"
	"os"
	"fmt"
	"im_go/im"
)

func StartClient(port int) {
	address := fmt.Sprintf("0.0.0.0:%d",port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
			auth.Token = "58b70f11-e280-40d8-8ddd-a3810bd65f7a"
			auth.DeviceId = "4c6aba79-f768-4e26-8344-aa2b7bc173ec"
			auth.PlatformType = 3
			p.Ver = 1
			p.Body = auth.ToData()
			//if err !=nil {
			//	fmt.Println("Marshal",err)
			//}
			p.Operation = OP_AUTH
			p.SeqId = 1
		}else if string(line) == "badauth"{
			var auth im.AuthenticationToken
			auth.Token = "fa9d0cdc-13f0-472c-99c2-e7b0100b8d09"
			auth.DeviceId = "4c6aba79-f768-4e26-8344-aa2b7bc173ec"
			auth.PlatformType = 3
			p.Ver = 1
			p.Body = auth.ToData()
			//if err !=nil {
			//	fmt.Println("Marshal",err)
			//}
			p.Operation = OP_AUTH
			p.SeqId = 1
		}else if string(line) == "badauth1"{
			var auth im.AuthenticationToken
			auth.Token = "badToken"
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
