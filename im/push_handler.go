package im

import (
	"im_go/model"
	"strconv"
	log "github.com/flywithbug/log4go"
)

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
}

