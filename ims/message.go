package ims

import (
	"fmt"
)

var message_descriptions map[int]string = make(map[int]string)

type Command int
func (cmd Command) String() string {
	c := int(cmd)
	if desc, ok := message_descriptions[c]; ok {
		return desc
	} else {
		return fmt.Sprintf("%d", c)
	}
}


type IMessage interface {
	ToData() []byte
	FromData(buff []byte) bool
}





