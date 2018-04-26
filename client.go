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
		}else {
			go send(conn,line)
		}
	}
}

func sendRandom(conn *net.TCPConn)  {
	send(conn,[]byte(RandString(1024)))
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