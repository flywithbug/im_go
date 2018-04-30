package ims

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"im_go/libs/proto"
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

type Message struct {
	proto.Proto
}

type IMessage interface {
	ToData() []byte
	FromData(buff []byte) bool
}

type AuthenticationToken struct {
	token       string
	platformId int8
	deviceId   string
}

func (auth *AuthenticationToken) ToData() []byte {
	var l int8
	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.BigEndian,auth.platformId)
	l = int8(len(auth.token))
	binary.Write(buffer,binary.BigEndian,l)
	buffer.Write([]byte(auth.token))
	l = int8(len(auth.deviceId))
	binary.Write(buffer,binary.BigEndian,l)
	buffer.Write([]byte(auth.deviceId))
	buf := buffer.Bytes()
	return buf
}

func (auth *AuthenticationToken) FromData(buff []byte) bool {
	var l int8
	if (len(buff)< 3) {
		return false
	}
	platformId := int8(buff[0])
	buffer := bytes.NewBuffer(buff[1:])
	binary.Read(buffer,binary.BigEndian,&l)
	if int(l) > buffer.Len() || int(l) < 0 {
		return false
	}
	token := make([]byte,l)
	buffer.Read(token)
	binary.Read(buffer,binary.BigEndian,&l)
	deviceId := make([]byte,l)
	buffer.Read(deviceId)
	auth.platformId = platformId
	auth.token = string(token)
	auth.deviceId = string(deviceId)
	return true
}


type AuthenticationStatus struct {
	status int32			//-1 权限校验未通过
	ip int32 //兼容版本0
}

func (auth *AuthenticationStatus) ToData(version int16) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, auth.status)
	if version == 0 {
		binary.Write(buffer, binary.BigEndian, auth.ip)
	}
	buf := buffer.Bytes()
	return buf
}

func (auth *AuthenticationStatus) FromData(version int, buff []byte) bool {
	if len(buff) < 4 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &auth.status)
	if version == 0 {
		if len(buff) < 8 {
			return false
		}
		binary.Read(buffer, binary.BigEndian, &auth.ip)
	}
	return true
}










