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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	in := bufio.NewReader(os.Stdin)
	go func() {
		for {
			if line, _, err := reader.ReadLine(); err == nil {
				log.Println(string(line),"------")

			}else {
				fmt.Println(err)
				break
			}
		}
	}()

	for {
		line, _, _ := in.ReadLine()
		// 模拟一个请求
		// {"command":"GET_CONN","data":null}
		// {"command":"GET_BUDDY_LIST","data":null}
		writer.WriteString(string(line) + "\n")
		writer.Flush()
	}

}
