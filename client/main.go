package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
)

func main()  {
	server := "127.0.0.1:5000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	in := bufio.NewReader(os.Stdin)
	for{
		line ,_,_ := in.ReadLine()
		conn.Write(line)
	}

}