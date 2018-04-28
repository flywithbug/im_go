package model

import (
	"time"
	"encoding/json"
)

/*
用户对象
*/
type User struct {
	Id       	string    `json:"id"`        //id
	Nick     	string    `json:"nick"`      //昵称
	Status   	string    `json:"status"`    //状态 0离线,1在线
	Sign     	string    `json:"sign"`      //个性签名
	Avatar   	string    `json:"avatar"`    //头像
	CreateAt 	time.Time `json:"create_at"` //注册日期
	UpdateAt 	time.Time `json:"update_at"` //更新日期
	Token    	string    `json:"token"`
	DeviceToken string	  `json:"device_token"` //设备推送token
	Platform	uint8	  `json:"platform"`		//设备类型 1 iOS，2 Android,3 Web
}

/*
 转JSON数据
*/
func (this *User) Encode() []byte {
	s, _ := json.Marshal(*this)
	return s
}

/*
 解析JSON
*/
func (this *User) Decode(data []byte) error {
	err := json.Unmarshal(data, this)
	return err
}
