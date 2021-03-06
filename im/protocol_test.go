package im

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriteHeader(t *testing.T) {

	//encode
	b := []byte("test")
	buffer := new(bytes.Buffer)
	var ph protoHeader
	ph.headerLen = RawHeaderSize
	ph.seq = 1
	ph.op = 2
	ph.bodyLen = int32(len(b))
	ph.ver = 1
	WriteHeader(ph, buffer)
	buffer.Write(b)
	bb := buffer.Bytes()
	fmt.Println(bb)

	//decode
	ph1, err := ReadHeader(bb)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ph1)
	if ph != *ph1 {
		fmt.Errorf("not expect")
	}

}
