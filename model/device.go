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
	Status 			int			`json:"status"`    //推送状态，默认1 推送，0 不推送
	Sound			int			`json:"sound"`     //推送是否有声音，1默认剩余，0 没有剩余，其他值待定
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

func (model *Device)SaveOrUpdateToDb()error  {
	//log.Info("%s",model.Environment)
	return SaveDeviceInfo(model.DeviceToken,model.DeviceId,model.UserAgent,model.UserId,model.UniqueMacUuid,model.Platform,model.Environment)
}

func UpdateDeviceInfo(deviceId string, status int)(int64,error)   {
	updateStmt, err := Database.Prepare("UPDATE im_device SET `status` = ? where device_id = ?")
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res, err := updateStmt.Exec(deviceId,status)
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"服务出错"}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return -1, &DatabaseError{"token已失效"}
	}
	return num,nil
}



func GetDevicesByUserId(userId string)([]Device,error)  {
	var devices  []Device
	rows, err := Database.Query("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment,status,Sound FROM im_device WHERE user_id = ?",userId)
	defer  rows.Close()
	if err != nil {
		log.Error(err.Error())
		return devices, &DatabaseError{"服务出错"}
	}
	for rows.Next()  {
		var device Device
		rows.Scan(&device.UserId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment,&device.Status,&device.Sound)
		devices = append(devices,device)
	}
	return devices,nil

}

func GetDevicesByDeviceId(deviceId string)(*Device,error)  {
	var device Device
	row := Database.QueryRow("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment,status ,Sound FROM im_device WHERE device_id = ?",deviceId)
	err := row.Scan(&device.UserId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment,&device.Status,&device.Sound)
		if err != nil {
			log.Error(err.Error()+deviceId)
			return nil, &DatabaseError{"未查询到对应的数据"}
		}
	return &device,nil
}


/*
 根据ID获取用户
//*/
//func GetUserByUId(uId string) (*User, error) {
//	var user User
//	row := Database.QueryRow("select id, app_id, user_id, nick, status, sign, avatar, create_at, update_at from im_user where id = ?", uId)
//	err := row.Scan(&user.Uid,&user.appId,&user.UserId, &user.Nick, &user.Status, &user.Sign, &user.Avatar, &user.createAt, &user.updateAt)
//	if err != nil {
//		log.Error(err.Error()+uId)
//		return nil, &DatabaseError{"根据ID查询用户-将结果映射至对象错误"}
//	}
//	return &user, err
//}