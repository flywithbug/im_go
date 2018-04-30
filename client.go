package main

import (
	"net"
	"fmt"
	"os"
	"im_go/imc"
	"bufio"
	"log"
	"im_go/libs/proto"
	"im_go/libs/bytes"
	"im_go/libs/define"
	gbytes "bytes"
	"encoding/binary"

	"github.com/pborman/uuid"
)


func main()  {

	tcpAddr, err := net.ResolveTCPAddr("tcp","0.0.0.0:9000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn ,err := net.DialTCP("tcp",nil,tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	c := imc.NewClient(conn)
	c.Listen()



	in := bufio.NewReader(os.Stdin)
	for  {
		line, _, _ := in.ReadLine()
		log.Println("readLine:",string(line))
		if string(line)== "auth" {
			AuthData(conn)
		}else {
			sendData(conn,line,define.OP_SEND_MSG)
		}
	}


}

func AuthData(conn *net.TCPConn)  {
	auth := new(authenticationToken2)
	auth.token = "c1428fb6-eb62-4360-ad01-ef6c84d7faa9"
	auth.deviceId = uuid.New()
	auth.platformId = 1
	fmt.Println("sendMSG",auth)
	sendData(conn,auth.ToData(),define.OP_AUTH)
}
func sendData(conn *net.TCPConn,data []byte,operation int32)  {
	fmt.Println("input:",string(data))
	p := new(proto.Proto)
	//p.Ver = 1
	p.Operation = operation
	//p.SeqId = int32(1)
	p.Body = []byte(data)
	//判断发送字符长度，过长提示
	wr := bufio.NewWriterSize(conn,len(p.Body)+50)
	b := bytes.NewWriterSize(len(p.Body)+50)
	p.WriteTo(b)
	_,err := wr.Write(b.Buffer())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(nn)
	wr.Flush()
}



//登录授权
type authenticationToken2 struct {
	token       string
	platformId int8
	deviceId   string
}



func (auth *authenticationToken2) ToData() []byte {
	var l int8
	buffer := new(gbytes.Buffer)
	binary.Write(buffer,binary.BigEndian,auth.platformId)
	l = int8(len(auth.token))
	binary.Write(buffer,binary.BigEndian,l)
	buffer.Write([]byte(auth.token))
	l = int8(len(auth.deviceId))
	binary.Write(buffer,binary.BigEndian,l)
	buffer.Write([]byte(auth.deviceId))
	buf := buffer.Bytes()
	return buf
}



func (auth *authenticationToken2) FromData(buff []byte) bool {
	var l int8
	if (len(buff)< 3) {
		return false
	}
	platformId := int8(buff[0])

	buffer := gbytes.NewBuffer(buff[1:])
	binary.Read(buffer,binary.BigEndian,&l)
	if int(l) > buffer.Len() || int(l) < 0 {
		return false
	}
	token := make([]byte,l)
	buffer.Read(token)


	binary.Read(buffer,binary.BigEndian,&l)
	deviceId := make([]byte,l)
	buffer.Read(deviceId)


	auth.platformId = platformId
	auth.token = string(token)
	auth.deviceId = string(deviceId)

	return true
}
