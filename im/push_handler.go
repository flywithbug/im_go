package im

import (
	"im_go/model"
	"strconv"
	log "github.com/flywithbug/log4go"
	"encoding/json"
	"im_go/http"
)

const POSTURLPATH  = "http://127.0.0.1:9012/api/push"

//type PushModel struct {
//	Title 			string 	`json:"title"`
//	SubTitle 		string	`json:"sub_title"`
//	BadgeNumber 	uint	`json:"badge_number"`
//	AppId       	int		`json:"app_id"` //用于索引证书
//	DeviceToken 	string	`json:"device_token"`
//	EnvironmentType	int		`json:"environment_type"` //默认0位production环境
//	Sound 			string	`json:"sound"`
//	Body            string	`json:"body"` //推送时显示的内容
//}
//
//type MessageBoddy struct {
//	Content   		string		`json:"content"`
//	Type 			int			`json:"type"`
//
//}


func PushServiceHandler(sender,receiver int32, appId int64,pro *Proto)  {
	user, err := model.GetUserByUId(strconv.Itoa(int(receiver)))
	if err != nil {
		log.Info(err.Error())
		return
	}
	sUser, err := model.GetUserByUId(strconv.Itoa(int(sender)))
	if err != nil {
		log.Info(err.Error())
		return
	}
	if user.GetAppId() != appId {
		log.Info("appId not equal")
		return
	}
	if pro.Operation == OP_MSG {
		device ,err := model.GetDeviceByUserId(user.UserId)
		if err != nil {
			log.Info(err.Error())
			return
		}
		msg := new(Message)
		msg.FromData(pro.Body)

		var msgBody = model.MessageBoddy{}
		err = json.Unmarshal(msg.body,&msgBody)
		if err != nil {
			log.Info(err.Error())
			return
		}
		push := model.PushModel{}
		push.DeviceToken = device.DeviceToken
		push.BadgeNumber ,err = model.MessageUnSendedCount(receiver)
		push.Title = sUser.Nick
		if len(msgBody.Content) > 0 {
			push.Body = msgBody.Content
		}else {
			push.Body = "新消息"
		}
		push.AppId = int(appId)
		if err != nil {
			log.Info(err.Error())
			return
		}
		push.EnvironmentType = device.Environment
		_ ,err = http.POST(POSTURLPATH,push,nil)
		if err != nil {
			log.Error(err.Error())
			return
		}
		//log.Info(msg.Description() + string(b))
	}


}

