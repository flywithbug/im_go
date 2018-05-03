package model

import (
	_ "database/sql"
	log "github.com/flywithbug/log4go"
	"errors"
)

type Device struct {
	Id  			int64   	`json:"id"`
	DeviceToken 	string	 	`json:"device_token"`
	DeviceId		string		`json:"device_id"`
	Platform		int			`json:"platform"`
	Description		string		`json:"description"`
}

func SaveDeviceInfo(token ,deviceId,description string,platformType int)error  {
	stmt,err :=Database.Prepare("INSERT into im_device SET device_id=? ON DUPLICATE key UPDATE ,device_token=?,platform=?,description=?")
	if err != nil{
		log.Warn(err.Error())
		err = errors.New("服务错误")
		return err
	}
	defer stmt.Close()
	_,err = stmt.Exec(deviceId,token,platformType,description)
	if err!= nil {
		log.Warn(err.Error())
		err = errors.New("device_id not found")
		return err
	}
	return nil
}

func GetDeviceBy(deviceId string)(*Device,error)  {
	var device Device
	row := Database.QueryRow("SELECT id,device_id,device_token,platform,description FROM im_device WHERE device_id = ?",deviceId)
	err := row.Scan(&device.Id,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.Description)
	if err != nil {
		return nil,&DatabaseError{"未查询到该设备"}
	}
	return &device,nil
}