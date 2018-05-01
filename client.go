package main


import (
	"bufio"
	"log"
	"net"
	"os"
	"io"
	"fmt"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9000")

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

func ioCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}