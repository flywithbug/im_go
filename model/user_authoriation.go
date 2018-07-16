package model

import _ "database/sql"
import log "github.com/flywithbug/log4go"


//CREATE TABLE `im_user_authorization` (
//`host_id` varchar(40) NOT NULL DEFAULT '',
//`guest_id` varchar(40) NOT NULL DEFAULT '',
//`status` tinyint(4) NOT NULL COMMENT '-1 拒绝，0 申请，1 接受',
//`type` tinyint(4) NOT NULL COMMENT '类型',
//`validity_time_stamp` varchar(12) NOT NULL DEFAULT '' COMMENT '有效期',
//PRIMARY KEY (`host_id`,`guest_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

type UserAuthorization struct {
	HostId 				string 		`json:"host_id"`
	GuestId				string		`json:"guest_id"`
	Status 				int			`json:"status"`
	AType        		int			`json:"a_type"`
}

func InsertAuthorization(hostId,guestId string,status,aType int)error  {
	insStmt,err := Database.Prepare("INSERT INTO im_user_authorization (host_id,guest_id,status,a_type) VALUES (?,?,?,?)")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer insStmt.Close()
	_, err = insStmt.Exec(hostId,guestId,status,aType)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetAuthorization(hostId,guestId string,aType int)(*UserAuthorization,error) {
	var authorization UserAuthorization
	row := Database.QueryRow("SELECT host_id,guest_id,statas, a_type  from im_user_authorization where host_id=? and guest_id=? AND a_type=?",hostId,guestId,aType)
	err := row.Scan(&authorization.HostId,&authorization.GuestId,&authorization.Status,&authorization.AType)
	if err != nil {
		log.Error(err.Error()+hostId+guestId)
		return nil, &DatabaseError{"未查询到该数据"}
	}
	return &authorization, nil
}

func UpdateAuthorization(hostId,guestId string,aType,status int)error  {
	_, err := GetAuthorization(hostId,guestId,aType)
	if err != nil {
		err = InsertAuthorization(hostId,guestId,status,aType)
		return err
	}
	updateStmt,err := Database.Prepare("UPDATE im_user_authorization SET `status` = ? WHERE  host_id=? and guest_id=? AND a_type=?")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res ,err := updateStmt.Exec(status,hostId,guestId,aType)
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return  &DatabaseError{"未查询到该用户"}
	}
	return nil


}
