package handler

import (
	"github.com/gin-gonic/gin"
	"im_go/model"
)

// 用户信息获取和填充
func User(ctx *gin.Context) (*model.User, bool) {
	o, ok := ctx.Get(KeyContextUser)
	if !ok {
		return nil, false
	}
	aUser, ok := o.(*model.User)
	if !ok {
		return nil, false
	}
	return aUser, true
}

func UserToken(ctx *gin.Context) (string, bool) {
	token, err := ctx.Cookie(KeyUserToken)
	if err != nil {
		return "", false
	}
	return token, true
}

func UserDeviceId(ctx *gin.Context) (string, bool)   {
	o, ok := ctx.Get(KeyContextDeviceId)
	if !ok {
		return "", false
	}
	deviceId, ok := o.(string)
	if !ok {
		return "", false
	}
	return deviceId, true
}