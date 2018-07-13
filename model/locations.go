package model

import _ "database/sql"

type Location struct {
	LType        int		`json:"l_type"`
	LTimeStamp 	string		`json:"l_time_stamp"`
	UserId 	 	string	   	`json:"user_id"`   //uuid生成
	Latitude    string		`json:"latitude"`   //维度
	Longitude   string		`json:"longitude"`   //经度
}



/*
 保存登录状态
 */
func SaveLocationsPath(userId, longitude,latitude ,time_stamp string,lType int) error {
	insStmt, err := Database.Prepare("INSERT into im_user_location_path (user_id, longitude,latitude ,l_time_stamp ,l_type) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return &DatabaseError{"服务错误:"+err.Error()}
	}
	defer insStmt.Close()
	_, err = insStmt.Exec(userId,longitude,latitude,time_stamp,lType)
	if err != nil {
		return &DatabaseError{"服务错误:"+err.Error()}
	}
	return nil
}







