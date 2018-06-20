package im

import (
	"im_go/model"
	"strconv"
	log "github.com/flywithbug/log4go"
	"im_go/http"
)

const POSTURLPATH  = "http://127.0.0.1:9012/api/push"

type PushModel struct {
	Title 			string 	`json:"title"`
	SubTitle 		string	`json:"sub_title"`
	BadgeNumber 	uint	`json:"badge_number"`
	AppId       	int		`json:"app_id"` //用于索引证书
	DeviceToken 	string	`json:"device_token"`
	EnvironmentType	int		`json:"environment_type"` //默认0位production环境
	Sound 			string	`json:"sound"`
	Body            string	`json:"body"` //推送时显示的内容
}


func PushServiceHandler(uId int32, appId int64,pro *Proto)  {
	user, err := model.GetUserByUId(strconv.Itoa(int(uId)))
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
		push := model.PushModel{}
		push.DeviceToken = device.DeviceToken
		push.BadgeNumber ,err = model.MessageUnSendedCount(uId)
		push.Body = "新消息"
		push.AppId = int(appId)
		if err != nil {
			log.Info(err.Error())
			return
		}
		push.EnvironmentType = device.Environment

		b ,err := http.POST(POSTURLPATH,push,nil)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info(string(b))
	}


}

