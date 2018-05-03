package model

import (
	_ "database/sql"
	log "github.com/flywithbug/log4go"

	"fmt"
)

type UserRelationShip struct {
	Id 			int			`json:"id"`
	UId 		int			`json:"u_id"`
	Status 		int			`json:"status"`
	RelationId 	string		`json:"relation_id"`
	FriendId	string		`json:"friend_id"`
	Remarks		string		`json:"remarks"`
}

func AddUserRelation(uId int,friendId int)(int64,error)  {
	inStmt,err := Database.Prepare("Replace INTO im_relationship SET relation_id = ?,u_id = ?,status = ?,friend_id = ?")
	defer inStmt.Close()
	if err != nil {
		log.Error(err.Error())
		return -1 ,&DatabaseError{"服务错误"}
	}
	relationId := fmt.Sprintf("%d-%d",uId,friendId)
	res,err := inStmt.Exec(relationId,uId,0,friendId)
	if err != nil {
		return -1 ,&DatabaseError{"服务错误"}
	}
	id,err := res.LastInsertId()
	if err != nil {
		return -1 ,&DatabaseError{"服务错误"}
	}
	return id,nil
}

func UpdateRelationRemark(relationId ,remark string)error  {
	updateStmt,err := Database.Prepare("UPDATE im_relationship SET `remark` = ? WHERE relation_id = ?")
	defer updateStmt.Close()
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务错误"}
	}
	_ ,err = updateStmt.Exec(remark,relationId)
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务错误"}
	}
	return nil
}

func DelRelationShip(relationId string)error  {
	delStmt ,err := Database.Prepare("DELETE FROM im_relationship WHERE relation_id = ?")
	defer delStmt.Close()
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务错误"}
	}
	_,err = delStmt.Exec(relationId)
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务错误"}
	}
	return nil
}


