package im

import (
	"testing"

	"fmt"
	"time"
)

func TestBaseMessage_Decode(t *testing.T) {
	msg := Message{
		receiver:  10001,
		sender:    10002,
		msgId:     20020,
		body:      []byte("hello world"),
		timestamp: time.Now().Unix(),
	}
	fmt.Println(msg.Description())

	b := msg.ToData()
	fmt.Println(msg.ToData())

	var m Message
	m.FromData(b)
	fmt.Println(m.Description())

}
