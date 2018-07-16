package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"im_go/model"
)

//CREATE TABLE `im_user_authorization` (
//`host_id` varchar(40) NOT NULL DEFAULT '',
//`guest_id` varchar(40) NOT NULL DEFAULT '',
//`status` tinyint(4) NOT NULL COMMENT '-1 拒绝，0 申请，1 接受',
//`type` tinyint(4) NOT NULL COMMENT '类型',
//`validity_time_stamp` varchar(12) NOT NULL DEFAULT '' COMMENT '有效期',
//PRIMARY KEY (`host_id`,`guest_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

func UpdateAuthorization(c *gin.Context) {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	auth := model.UserAuthorization{}
	err := c.BindJSON(&auth)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}
	err = model.UpdateAuthorization(auth.HostId,auth.GuestId,auth.AType,auth.Status)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError, "server error"+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}

