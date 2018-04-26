package main

import (
	"net"
	"fmt"
	"os"

	"math/rand"

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

	defer conn.Close()

	//in := gbufio.NewReader(os.Stdin)
	//
	//
	//buffer := make([]byte, 1024)
	p := new(proto.Proto)
	p.Ver = 1
	p.Operation = define.OP_AUTH
	p.SeqId = int32(0)
	p.Body = []byte("test")
	wr := bufio.NewWriterSize(conn,len(p.Body)+50)
	b := bytes.NewWriterSize(1024*10)
	p.WriteTo(b)
	wr.Write(b.Buffer())
	//for {
	//	line, _, _ := in.ReadLine()
	//	fmt.Println(line)
	//	conn.Write(line)
	//	readLen, _ := conn.Read(buffer)
	//	fmt.Println(readLen,string(buffer))
	//}
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