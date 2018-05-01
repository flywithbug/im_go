package im

import (
	log "github.com/golang/glog"
	"fmt"
	"bytes"
	"encoding/binary"
	"im_go/model"
	"encoding/json"
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
	fmt.Println("authToken:",auth.Token,"platform:",auth.PlatformType,"deviceId:",auth.DeviceId)


	login,err := model.GetLoginByToken(auth.Token,model.STATUS_LOGIN)
	if err != nil {
		fmt.Println(err)
		var authStatus  AuthenticationStatus
		authStatus.Status= -1
		pro.Operation = OP_AUTH_REPLY
		pro.Body = authStatus.ToData()
		pro.SeqId = 0
		client.EnqueueMessage(pro)
		return
	}
	jb,_ :=json.Marshal(login)
	fmt.Println(string(jb))


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




type AuthenticationStatus struct {
	Status int32   //-1 验证失败
}


func (auth *AuthenticationStatus) ToData ()[]byte {
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