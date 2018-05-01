package im

import (
	"net"
	"log"
	"bufio"
)

type Client struct {
	conn *net.TCPConn
	reader *bufio.Reader //读取
	writer *bufio.Writer //输出

	Out    chan  []byte   //输出消息


}

func NewClient(conn *net.TCPConn)*Client  {
	client := new(Client)
	client.conn = conn
	client.Out = make(chan []byte,10)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	client.reader =reader
	client.writer = writer

	return client
}


func (client *Client)Read()  {
	for{
		//log.Printf("Read")
		if line, _,err := client.reader.ReadLine();err == nil{
			log.Println("read",string(line))
			go func() {
				client.Out <- line
			}()
		}else {
			return
		}
	}
}

func (client*Client)SendMsg(in string)  {
	client.conn.Write([]byte(in +"\n"))
}



func (client *Client)Write()  {
	for b := range client.Out{

		client.SendMsg(string(b))
	}
}


func (client *Client) Listen() {
	go client.Read()
	go client.Write()

}
