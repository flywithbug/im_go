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
	Status 			bool		`json:"status"`    //推送状态，默认1 推送，0 不推送
	Sound			int			`json:"sound"`     //推送是否有声音，1默认声音，0 没有声音，其他值待定
	ShowDetail		bool		`json:"show_detail"` //是否显示推送消息详情
}


func (model *Device)SaveOrUpdateToDb()error  {
	dev,_ := GetDevicesByDeviceId(model.DeviceId)
	if dev != nil {
		_,err := model.UpdateDeviceInfo()
		if err != nil {
			log.Info(err.Error())
		}
		return err
	}
	return SaveDeviceInfo(model.DeviceToken,model.DeviceId,model.Platform,model.UserAgent,model.UserId,model.UniqueMacUuid,model.Environment,model.Status,model.Sound,model.ShowDetail)
}

func SaveDeviceInfo(device_token,device_id string,platform int,user_agent,user_id,unique_mac_uuid string, environment int,status bool,sound int, show_detail bool)error  {
	stmt,err :=Database.Prepare("INSERT into im_device  (device_token,device_id,platform,user_agent,user_id,unique_mac_uuid,environment,status,sound , show_detail ) VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil{
		log.Warn(err.Error())
		err = errors.New("服务错误")
		return err
	}
	defer stmt.Close()
	_,err = stmt.Exec(device_token,device_id,platform,user_agent,user_id,unique_mac_uuid,environment,status,sound,show_detail)
	if err!= nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}


func (devicd *Device)UpdateDeviceInfo()(int64,error)  {
	updateStmt, err := Database.Prepare("UPDATE im_device SET user_id=? ,device_token=?,platform=?,user_agent=? ,unique_mac_uuid = ?,environment= ?,status = ? ,sound =?, show_detail=? where device_id = ?")
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res, err := updateStmt.Exec(devicd.UserId,devicd.DeviceToken,devicd.Platform,devicd.UserAgent,devicd.UniqueMacUuid,devicd.Environment,devicd.Status,devicd.Sound,devicd.ShowDetail,devicd.DeviceId)
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"服务出错"}
	}
	num, err := res.RowsAffected()
	if err != nil {
		log.Info(err.Error())
		return -1, &DatabaseError{"更新失败"}
	}
	return num,nil

}




func GetDevicesByUserId(userId string)([]Device,error)  {
	var devices  []Device
	rows, err := Database.Query("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment,status,sound, show_detailFROM im_device WHERE user_id = ?",userId)
	defer  rows.Close()
	if err != nil {
		log.Error(err.Error())
		return devices, &DatabaseError{"服务出错"}
	}
	for rows.Next()  {
		var device Device
		rows.Scan(&device.UserId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment,&device.Status,&device.Sound,&device.ShowDetail)
		devices = append(devices,device)
	}
	return devices,nil

}

func GetDevicesByDeviceId(deviceId string)(*Device,error)  {
	var device Device
	row := Database.QueryRow("SELECT user_id,device_id,device_token,platform,user_agent,unique_mac_uuid,environment,status ,sound ,show_detail FROM im_device WHERE device_id = ?",deviceId)
	err := row.Scan(&device.UserId,&device.DeviceId,&device.DeviceToken,&device.Platform,&device.UserAgent,&device.UniqueMacUuid,&device.Environment,&device.Status,&device.Sound,&device.ShowDetail)
	if err != nil {
		log.Error(err.Error()+deviceId)
		return nil, &DatabaseError{"未查询到对应的数据"}
	}
	return &device,nil
}

