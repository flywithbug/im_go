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
	Platform		int			`json:"platform"`   //1.iOS, 2.android
	UserAgent		string		`json:"user_agent"`
	UserId 			string		`json:"user_id"`  //绑定的用户UserId
	UniqueMacUuid	string		`json:"unique_mac_uuid"`
	Environment     int			`json:"environment"`  //客户端开发环境 默认production:0, development:1 环境
}

func SaveDeviceInfo(deviceToken ,deviceId,user_agent ,userId,unique_mac_uuid string,platformType,environment int)error  {
	stmt,err :=Database.Prepare("INSERT into im_device SET user_id=? ,device_id=?,device_token=?,platform=?,user_agent=? ,unique_mac_uuid = ?,environment= ? ON DUPLICATE key UPDATE device_token=?,user_id=?,platform=?,user_agent=?,unique_mac_uuid = ?,environment=? ")
	if err != nil{
		log.Warn(err.Error())
		err = errors.New("服务错误")
		return err
	}
	defer stmt.Close()
	_,err = stmt.Exec(userId,deviceId,deviceToken,platformType,user_agent,unique_mac_uuid,environment,deviceToken,userId,platformType,user_agent,unique_mac_uuid,environment)
	if err!= nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (model *Device)SaveToDb()error  {
	//log.Info("%s",model.Environment)
	return SaveDeviceInfo(model.DeviceToken,model.DeviceId,model.UserAgent,model.UserId,model.UniqueMacUuid,model.Platform,model.Environment)
}

func GetDeviceByUserId(userId string)(*Device,error)  {
	//log.Info(userId)
	var device Device
	row := Database.QueryRow("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment FROM im_device WHERE user_id = ?",userId)
	err := row.Scan(&userId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment)
	if err != nil {
		return nil,&DatabaseError{"未查询到该设备"}
	}
	return &device,nil
}

func GetDevicesByUserId(userId string)([]Device,error)  {
	var devices  []Device
	rows, err := Database.Query("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment FROM im_device WHERE user_id = ?",userId)
	defer  rows.Close()
	if err != nil {
		log.Error(err.Error())
		return devices, &DatabaseError{"服务出错"}
	}
	for rows.Next()  {
		var device Device
		rows.Scan(&userId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment)
		devices = append(devices,device)
	}
	return devices,nil

}



////查找发送人消息
//func FindeMessagesSender(sender int32,status int)([]IMMessage,error)  {
//	var messages []IMMessage
//	rows ,err := Database.Query("SELECT id,sender,receiver,content,time_stamp,status,update_at FROM im_message WHERE sender = ? AND status = ?",sender,status)
//	defer rows.Close()
//	if err != nil {
//		log.Error(err.Error())
//		return messages, &DatabaseError{"服务出错"}
//	}
//	for rows.Next(){
//		var msg IMMessage
//		rows.Scan(&msg.Id,&msg.Sender,&msg.Receiver,&msg.Content,&msg.TimeStamp,&msg.Status,&msg.UpdateAt)
//		messages = append(messages, msg)
//	}
//	return messages,nil
//}