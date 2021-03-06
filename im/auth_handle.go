package im

import (
	"fmt"
	"strings"
	"time"

	log "github.com/flywithbug/log4go"

	"im_go/model"
	"sync/atomic"
)

const (
	AuthenticationStatusSuccess  = 1
	AuthenticationStatusBadToken = -1
	AuthenticationStatusBadLogin = -2
)

func (client *Client) HandleAuthToken(pro *Proto) {
	var auth AuthenticationToken
	if !auth.FromData(pro.Body) {
		log.Info("AuthenticationToken decode error, auth:%s",auth.Description())
		return
	}

	if client.uid > 0 && strings.EqualFold(client.Token, auth.Token) {
		log.Info("repeat login")
		return
	}

	//call back Body
	var authStatus AuthenticationStatus
	pro.Operation = OP_AUTH_ACK
	login, err := model.GetLoginByToken(auth.Token)
	if err != nil {
		log.Warn("get login by token error, auth:%s",auth.Description())
		log.Warn(err.Error())
		authStatus.Status = AuthenticationStatusBadToken
		pro.Body = authStatus.ToData()
		client.EnqueueMessage(*pro)
		return
	}

	if login.Status != 1 {
		log.Error("token 已失效")
		authStatus.Status = AuthenticationStatusBadToken
		pro.Body = authStatus.ToData()
		client.EnqueueMessage(*pro)
		return
	}

	if login != nil && (login.UId == 0 || login.AppId == 0 || login.UserId == "") {
		errString := fmt.Sprintf("auth Error uid:%d appId:%d userId:%s", login.UId, login.AppId, login.UserId)
		log.Error(errString)
		authStatus.Status = AuthenticationStatusBadLogin //登录用户信息有误
		pro.Body = authStatus.ToData()
		client.EnqueueMessage(*pro)
		return
	}

	//发消息给其他客户端登录的用户下线，并关闭其他客户端的connection
	//暂时只能单端登录
	authStatus.Status = AuthenticationStatusSuccess
	pro.Body = authStatus.ToData()
	send := client.EnqueueMessage(*pro)

	if send {
		client.version = pro.Ver

		client.appId = login.AppId
		client.uid = login.UId
		client.userId = login.UserId
		client.platformId = auth.PlatformType
		client.Token = auth.Token
		client.online = true
		client.forbidden = login.Forbidden
		client.tm = time.Now()

		//clientInfo := fmt.Sprintf("auth token:%s appid:%d uid:%d device id:%s forbidden:%d",
		//	login.Token, client.appid, client.uid, client.deviceId, client.forbidden)
		//log.Debug(clientInfo)

		client.AddClient()
		atomic.AddInt64(&serverSummary.nclients, 1)

		//用户授权成功之后 发送离线消息
		client.sendOffLineMessage()
		//登出其他账号
		//client.LogOutOtherClient()
		//model.UpdateUserStatus(login.UId, model.STATUS_LOGIN)
	}else {
		log.Error("auth status  消息返回客户端失败")
	}
}



func (client *Client) LogOutOtherClient() {
	//p := new(Proto)
	//p.Operation = OP_DISCONNECT_ACK
	//route := appRoute.FindRoute(client.appid)
	//clients := route.FindClientSet(client.uid)
	////可以扩展多端同时在线。
	//for c, _ := range clients {
	//	//不再发送给自己
	//	if c == client {
	//		continue
	//	}
	//	//发送踢出消息
	//	c.EnqueueMessage(p)
	//	//c.handleClientClosed()
	//}
	//model.LogoutOthers(client.Token, client.uid)
}

