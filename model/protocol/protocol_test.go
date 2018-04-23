package protocol

import (
	"testing"
	"fmt"
)



func TestPacket(t *testing.T) {
	//matchString := "www.flywithbug.com:90089-----flywithbug-----"
	//m := []byte(matchString)
	//match := lengthtRe.FindSubmatch(m)
	//for _,m := range match{
    	//fmt.Println(string(m))
	//}

	content := "2322323232323232323"
	msg := Packet([]byte(content))
	fmt.Println("msg:",string(msg))

	headers ,err := GetHeader(msg)
	if err != nil {
		panic(err)
		return
	}
	for _,k := range headers{
		fmt.Println("match:",string(k))
	}
	fmt.Println("headerlen",len(headers[0]))
}