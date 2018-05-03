package main

import (
	"flag"
	"im_go/imc"

)

func main() {
	flag.Parse()
	imc.StartClient(9000)
}
