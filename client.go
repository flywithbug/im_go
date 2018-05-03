package main

import (
	"flag"
	"im_go/imc"
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Unix())
	flag.Parse()
	imc.StartClient(9000)
}
