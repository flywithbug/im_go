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
}

func SaveDeviceInfo(deviceToken ,deviceId,description string,platformType int,userId string)error  {
	stmt,err :=Database.Prepare("INSERT into im_device SET user_id=? ,device_id=?,device_token=?,platform=?,description=? ON DUPLICATE key UPDATE device_id=?,device_token=?,platform=?,description=?")
	if err != nil{
		log.Warn(err.Error())
		err = errors.New("服务错误")
		return err
	}
	fmt.Println()
	defer stmt.Close()
	_,err = stmt.Exec(userId,deviceId,deviceToken,platformType,description,deviceId,deviceToken,platformType,description)
	if err!= nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (model *Device)SaveToDb()error  {
	return SaveDeviceInfo(model.DeviceToken,model.DeviceId,model.Description,model.Platform,model.UserId)
}

func GetDeviceByUserId(userId string)(*Device,error)  {
	var device Device
	row := Database.QueryRow("SELECT id,user_id,device_id,device_token,platform,description FROM im_device WHERE user_id = ?",userId)
	err := row.Scan(&userId,&device.Id,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.Description)
	if err != nil {
		return nil,&DatabaseError{"未查询到该设备"}
	}
	return &device,nil
}