package main

import (
	"fmt"
	"im_go/ims"
	"im_go/libs/perf"
)

func main()  {
	fmt.Printf("hello world")
	perf.Init([]string{"localhost:6971"})

	ims.Listen(9000)
}
