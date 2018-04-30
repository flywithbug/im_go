package imc

import (
	"net"
	"im_go/libs/bufio"
	"im_go/libs/proto"
	"fmt"
	"io"
)

type Client struct {
	conn 	*net.TCPConn
	reader 	*bufio.Reader	//读取
	writer 	*bufio.Writer	//输出

	in  	chan  	*proto.Proto
	msg 	*proto.Proto
	publicIp int32

}

func NewClient(conn *net.TCPConn)(client *Client)  {
	client = new(Client)
	client.conn = conn
	addr := conn.RemoteAddr()
	if ad, ok := addr.(*net.TCPAddr); ok {
		fmt.Println("address",ad)
		ip4 := ad.IP.To4()
		if len(ip4)>4 {
			client.publicIp = int32(ip4[0]) << 24 | int32(ip4[1]) << 16 | int32(ip4[2]) << 8 | int32(ip4[3])
		}
	}
	client.msg = new(proto.Proto)
	client.in = make(chan *proto.Proto,10)
	reader := bufio.NewReaderSize(conn,int(proto.MaxPackSize))
	writer := bufio.NewWriter(conn)
	client.reader =reader
	client.writer = writer
	return client
}

func (client *Client)Listen()  {
	go client.Read()
	go client.Write()
}

func (client *Client)read()(*proto.Proto,error)  {
	err := client.msg.ReadTCP(client.reader)
	return client.msg,err
}

func (client *Client)send(msg *proto.Proto)error  {
	return msg.WriteTCP(client.writer)
}

func (client *Client)Read()  {
	for{

		if msg,err := client.read();err != nil{
			if err == io.EOF {
				fmt.Println("close:",err)
				break
			}
			fmt.Println("else:",err)

		}else {
			fmt.Println("receiveMsg",msg)
		}
	}
}

func (client *Client)Write()  {
	for {
		select {
		case msg := <- client.in:
			fmt.Println(msg)
			err := client.send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}