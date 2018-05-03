package imc

import (
	"bufio"
	"net"
	"os"
	"fmt"
	"im_go/im"
	"time"
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
			case OP_AUTH_ACK:
				var auth im.AuthenticationStatus
				auth.FromData(msg.Body)
				fmt.Println("授权状态",auth.Status)
			case OP_SEND_MSG_ACK:
				var ack MessageACK
				ack.FromData(msg.Body)
				fmt.Println(ack.Description(),msg.Description())
			}
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
		}else if string(line) == "msg"{
			msg := Message{
				receiver:10002,
				sender:10001,
				msgId:2323232,
				body:[]byte("msg"),
				timestamp:time.Now().Unix(),
			}
			p.Ver = 1
			p.Body = msg.ToData()

			p.Operation = OP_SEND_MSG
			p.SeqId = 123456
		}else if string(line) == "auth1"{
			var auth im.AuthenticationToken
			auth.Token = "2c06eaf6-e14a-4d06-ba42-15de3f11741a"
			auth.DeviceId = "4c6aba79-f768-4e26-8344-aa2b7bc173ec"
			auth.PlatformType = 1
			p.Ver = 1
			p.Body = auth.ToData()
			//if err !=nil {
			//	fmt.Println("Marshal",err)
			//}
			p.Operation = OP_AUTH
			p.SeqId = 1
		}else if string(line) == "msg1"{
			msg := Message{
				receiver:10001,
				sender:10002,
				msgId:20020,
				body:[]byte("msg1"),
				timestamp:time.Now().Unix(),
			}
			p.Ver = 1
			p.Operation = OP_SEND_MSG
			p.SeqId = 1
			p.Body = msg.ToData()
		}else if string(line) == "auth3"{ //另一个账号
			var auth im.AuthenticationToken
			auth.Token = "2c06eaf6-e14a-4d06-ba42-15de3f11741a"
			auth.DeviceId = "4c6aba79-f768-4e26-8344-aa2b7bc173ec"
			auth.PlatformType = 3
			p.Ver = 1
			p.Body = auth.ToData()
			//if err !=nil {
			//	fmt.Println("Marshal",err)
			//}
			p.Operation = OP_AUTH
			p.SeqId = 1
		}else if string(line) == "auth4"{
			var auth im.AuthenticationToken
			auth.Token = "3680e5ef-d5e8-4efe-82c0-9ad03cf16025"
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
