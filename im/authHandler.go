package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	log "github.com/flywithbug/log4go"

	"im_go/model"
	"sync/atomic"
)

const (
	AuthenticationStatusSuccess  = 0
	AuthenticationStatusBadToken = -1
	AuthenticationStatusBadLogin = -2
	AuthenticationStatusKickOut = -3
)

type AuthenticationToken struct {
	Token        string `json:"token"`
	PlatformType int8   `json:"platform_type"`
	DeviceId     string `json:"device_id"`
}


func (client *Client) HandleAuthToken(pro *Proto) {
	var auth AuthenticationToken
	auth.FromData(pro.Body)
	logInfo := fmt.Sprintf("authToken:%s platform:%d deviceId:%s", auth.Token, auth.PlatformType, auth.DeviceId)
	log.Info(logInfo)
	//call back Body

	if client.uid > 0 && strings.EqualFold(client.Token, auth.Token) {
		log.Info("repeat login")
		return
	}
	var authStatus AuthenticationStatus
	pro.Operation = OP_AUTH_REPLY
	pro.SeqId = 0

	login, err := model.GetLoginByToken(auth.Token)
	if err != nil {
		log.Error(err.Error())
		authStatus.Status = AuthenticationStatusBadToken
		pro.Body = authStatus.ToData()
		client.EnqueueMessage(pro)
		return
	}

	if login.Status != 1 {
		log.Error("token 已失效")
		authStatus.Status = AuthenticationStatusBadToken
		pro.Body = authStatus.ToData()
		client.EnqueueMessage(pro)
		return
	}

	if login != nil && (login.UId == 0 || login.AppId == 0 || login.UserId == "") {
		errString := fmt.Sprintf("auth Error uid:%d appId:%d userId:%s", login.UId, login.AppId, login.UserId)
		log.Error(errString)
		authStatus.Status = AuthenticationStatusBadLogin //登录用户信息有误
		pro.Body = authStatus.ToData()
		client.EnqueueMessage(pro)
		return
	}
	//发消息给其他客户端登录的用户下线，并关闭其他客户端的connection
	client.appid = login.AppId
	client.uid = login.UId
	client.userId = login.UserId
	client.platformId = auth.PlatformType
	client.Token = auth.Token
	client.online = true
	client.forbidden = login.Forbidden
	client.tm = time.Now()

	clientInfo := fmt.Sprintf("auth token:%s appid:%d uid:%d device id:%s forbidden:%d",
		login.Token, client.appid, client.uid, client.deviceId, client.forbidden)
	log.Info(clientInfo)

	authStatus.Status = AuthenticationStatusSuccess
	pro.Body = authStatus.ToData()
	client.EnqueueMessage(pro)

	client.AddClient()
	atomic.AddInt64(&serverSummary.nclients,1)
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
