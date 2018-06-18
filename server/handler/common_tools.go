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