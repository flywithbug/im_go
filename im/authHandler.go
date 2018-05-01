package im

import (
	log "github.com/golang/glog"
	"fmt"
	"bytes"
	"encoding/binary"
)


type AuthenticationToken struct {
	Token       	string		`json:"token"`
	PlatformType 	int8		`json:"platform_type"`
	DeviceId   		string		`json:"device_id"`
}



func (client *Client)HandleAuthToken(pro *Proto)  {
	if client.uid >0 {
		log.Info("repeat login")
		return
	}
	var auth AuthenticationToken
	auth.FromData(pro.Body)
	fmt.Println("authToken",auth.Token,auth.PlatformType,auth.DeviceId)
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
	if (len(buff) <= 3) {
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
