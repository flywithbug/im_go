package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//首次连接token
type AuthenticationToken struct {
	PlatformType int8   `json:"platform_type"`
	Token        string `json:"token"`
	DeviceId     string `json:"device_id"`
}

func (auth *AuthenticationToken) ToData() []byte {
	var l int8

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, auth.PlatformType)

	l = int8(len(auth.Token))
	binary.Write(buffer, binary.BigEndian, l)
	buffer.Write([]byte(auth.Token))

	l = int8(len(auth.DeviceId))
	binary.Write(buffer, binary.BigEndian, l)
	buffer.Write([]byte(auth.DeviceId))

	buf := buffer.Bytes()
	return buf
}

func (auth *AuthenticationToken) FromData(buff []byte) bool {
	var l int8
	if len(buff) <= 3 {
		return false
	}
	platformType := int8(buff[0])

	buffer := bytes.NewBuffer(buff[1:])

	binary.Read(buffer, binary.BigEndian, &l)
	if int(l) > buffer.Len() || int(l) < 0 {
		return false
	}
	token := make([]byte, l)
	buffer.Read(token)

	binary.Read(buffer, binary.BigEndian, &l)
	if int(l) > buffer.Len() || int(l) < 0 {
		return false
	}
	deviceId := make([]byte, l)
	buffer.Read(deviceId)

	auth.PlatformType = platformType
	auth.Token = string(token)
	auth.DeviceId = string(deviceId)
	return true
}

func (auth *AuthenticationToken) Description() string {
	return fmt.Sprintf("Token:%s,PlatformType:%d,DeviceId:%s" , auth.Token, auth.PlatformType,auth.DeviceId)
}



//授权状态ack
type AuthenticationStatus struct {
	Status int32 //-1 验证失败  -2 用户信息有误
}

func (auth *AuthenticationStatus) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, auth.Status)
	buf := buffer.Bytes()
	return buf
}

func (auth *AuthenticationStatus) FromData(buff []byte) bool {
	if len(buff) < 4 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &auth.Status)

	return true
}

