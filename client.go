package main

import (
	"im_go/imc"
	"flag"
)

func main() {
	flag.Parse()
	imc.StartClient(9000)
}

