package main

import (
	"net"
	"fmt"
	"os"
	"math/rand"
	gbufio "bufio"
	"strconv"
	"strings"
	"time"
	"im_go/libs/proto"
	"im_go/libs/define"
	"im_go/libs/bytes"
	"im_go/libs/bufio"
	"encoding/binary"
	gbytes "bytes"

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
	in := gbufio.NewReader(os.Stdin)
	for  {
		line, _, _ := in.ReadLine()
		fmt.Println("line",string(line))
		if string(line) == "random" {
			go sendRandom(conn)
		}else if string(line) == "auth"{
			Auth(conn)
		} else {
			go send(conn,line)
		}
	}


}

func sendRandom(conn *net.TCPConn)  {
	string := RandString(1024*10)
	count := 0
	for {
		go send(conn,[]byte(string))
		count++
		if count == 1000{
			break
		}
	}
}




func send(conn *net.TCPConn,data []byte)  {
	fmt.Println("input:",string(data))
	p := new(proto.Proto)
	p.Ver = 1
	p.Operation = define.OP_AUTH
	p.SeqId = int32(0)
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




func Auth(conn *net.TCPConn)  {
	auth := new(authenticationToken)
	auth.token = "token"
	auth.deviceId = "deviceId"
	auth.platformId = 1

	send(conn,auth.ToData())
}

/**
*生成随机字符
**/
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}




//登录授权
type authenticationToken struct {
	token       string
	platformId int8
	deviceId   string
}



func (auth *authenticationToken) ToData() []byte {
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



func (auth *authenticationToken) FromData(buff []byte) bool {
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
