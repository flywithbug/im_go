package model

import (
	_ "database/sql"
	log "github.com/flywithbug/log4go"
	"errors"
	"fmt"
)

type Device struct {
	Id  			int64   	`json:"id"`
	DeviceToken 	string	 	`json:"device_token"`
	DeviceId		string		`json:"device_id"`
	Platform		int			`json:"platform"`   //1.iOS, 2.android
	UserAgent		string		`json:"user_agent"`
	UserId 			string		`json:"user_id"`  //绑定的用户UserId
	UniqueMacUuid	string		`json:"unique_mac_uuid"`
	Environment     int			`json:"environment"`  //客户端开发环境 默认production:0, development:1 环境
}

func SaveDeviceInfo(deviceToken ,deviceId,user_agent string,platformType,environment int,userId,unique_mac_uuid string)error  {
	stmt,err :=Database.Prepare("INSERT into im_device SET user_id=? ,device_id=?,device_token=?,platform=?,user_agent=? ,unique_mac_uuid = ?,environment= ? ON DUPLICATE key UPDATE device_id=?,device_token=?,platform=?,user_agent=?,unique_mac_uuid = ?,environment=? ")
	if err != nil{
		log.Warn(err.Error())
		err = errors.New("服务错误")
		return err
	}
	fmt.Println()
	defer stmt.Close()
	_,err = stmt.Exec(userId,deviceId,deviceToken,platformType,user_agent,unique_mac_uuid,environment,deviceId,deviceToken,platformType,user_agent,unique_mac_uuid,environment)
	if err!= nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (model *Device)SaveToDb()error  {
	return SaveDeviceInfo(model.DeviceToken,model.DeviceId,model.UserAgent,model.Platform,model.Environment,model.UserId,model.UniqueMacUuid)
}

func GetDeviceByUserId(userId string)(*Device,error)  {
	log.Info(userId)
	var device Device
	row := Database.QueryRow("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment FROM im_device WHERE user_id = ?",userId)
	err := row.Scan(&userId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment)
	if err != nil {
		return nil,&DatabaseError{"未查询到该设备"}
	}
	return &device,nil
}