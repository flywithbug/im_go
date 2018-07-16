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
type paraAuthorizationM struct {
	UserId 				string 		`json:"user_id"`
	AType        		int			`json:"a_type"`
	Status 				int			`json:"status"`
	AsHost              bool		`json:"as_host"`
}


func UpdateAuthorization(c *gin.Context) {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	auth := paraAuthorizationM{}
	err := c.BindJSON(&auth)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}
	user , _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"no user")
		return
	}
	err = model.UpdateAuthorization(user.UserId,auth.UserId,auth.AType,auth.Status)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError, "server error"+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK,"success")
}

func GetAuthorizationStatus(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	para := paraAuthorizationM{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid "+err.Error())
		return
	}
	if len(para.UserId) == 0{
		aRes.SetErrorInfo(http.StatusBadRequest, "userId can not be nil")
		return
	}
	user , _ := User(c)
	if user == nil {
		aRes.SetErrorInfo(http.StatusBadRequest ,"no user")
		return
	}

	if para.AsHost{
		auth,err := model.GetAuthorization(user.UserId,para.UserId,para.AType)
		if err != nil {
			aRes.SetErrorInfo(http.StatusInternalServerError ,"server error 未授权"+err.Error())
			return
		}
		aRes.AddResponseInfo("auth",auth)
	}else {
		auth,err := model.GetAuthorization(para.UserId,user.UserId,para.AType)
		if err != nil {
			aRes.SetErrorInfo(http.StatusInternalServerError ,"server error 未授权"+err.Error())
			return
		}
		aRes.AddResponseInfo("auth",auth)
	}


}

func GetUserCurrentLocations(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK, aRes)
	}()
	para := paraAuthorizationM{}
	err := c.BindJSON(&para)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid "+err.Error())
		return
	}
	if len(para.UserId) == 0{
		aRes.SetErrorInfo(http.StatusBadRequest, "userId can not be nil")
		return
	}
	user, err := model.GetUserWithLocationByUserId(para.UserId)
	if err != nil {
		aRes.SetErrorInfo(http.StatusInternalServerError, "server error"+err.Error())
		return
	}
	aRes.AddResponseInfo("user",user)
}



