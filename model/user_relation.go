package model

import (
	_ "database/sql"

	log "github.com/flywithbug/log4go"
	"strings"
)

type UserRelationShip struct {
	UserId   	string			`json:"user_id"`
	FUserId  	string			`json:"f_user_id"`
	ReceiverId  string			`json:"receiver_id"`
	Status 		int				`json:"status"`  //-2拉黑，-1 拒绝，0 申请，1 接受
}


/*
 申请好友
*/
func ApplyUserRelationShip(userId,friend_userId,receiverId string) error {
	if userId > friend_userId {
		temp := friend_userId
		friend_userId = userId
		userId = temp
	}

	uRelation,_ := GetUserRelationByUserIdAndFriendId(userId,friend_userId)
	if uRelation != nil {
		if uRelation.Status < 1 && uRelation.Status> -3 {
			err := UpdateUserRelation(userId,friend_userId,0)		
			return  err
		}else if uRelation.Status == 1 {
			return &DatabaseError{"已经是好友"}
		}else {
			return &DatabaseError{"您已被拉黑，无法添加好友"}
		}
	}
	insStmt, err := Database.Prepare("INSERT INTO im_user_relation (f_user_id,user_id,receiver_id,status) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer insStmt.Close()
	_, err = insStmt.Exec(friend_userId,userId,receiverId,0)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetUserRelationByUserIdAndFriendId(userId,friend_userId string)(*UserRelationShip,error)  {
	var relation UserRelationShip
	row := Database.QueryRow("SELECT f_user_id,user_id,receiver_id,status  from im_user_relation where user_id=? and f_user_id=? ",userId,friend_userId)
	err := row.Scan(&relation.FUserId,&relation.UserId,&relation.ReceiverId,&relation.Status)
	if err != nil {
		log.Error(err.Error()+userId)
		return nil, &DatabaseError{"根据账号及密码查询用户错误"}
	}
	return &relation, nil
}


//更新好友关系
func UpdateUserRelation(userId,friendId string,status int)error  {
	if userId > friendId {
		temp := friendId
		friendId = userId
		userId = temp
	}
	updateStmt,err := Database.Prepare("UPDATE im_user SET `status`= ?  WHERE user_id=? AND  f_user_id= ?")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res ,err := updateStmt.Exec(status,userId,friendId)
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return  &DatabaseError{"未查询到该用户"}
	}
	return nil
}

func GetFriendRelations(userId string)([]UserRelationShip,error)  {
	var relations []UserRelationShip
	rows , err := Database.Query("SELECT user_id, f_user_id from im_user_relation WHERE  user_id = ? OR f_user_id=?",userId,userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var re UserRelationShip
		rows.Scan(&re.UserId,&re.FUserId)
		relations = append(relations, re)
	}
	return relations,nil
}

func GetFriendsList(userId string)([]SimpleUser,error)  {
	relations ,err := GetFriendRelations(userId)
	if err != nil {
		return nil, err
	}
	var users []SimpleUser
	for _, re := range relations {
		tempId := re.UserId
		if strings.EqualFold(re.UserId,userId) {
			tempId = re.FUserId
		}
		su,_ := GetUserWithLocationByUserId(tempId)
		if err == nil {
			users = append(users, *su)
		}
	}
	return users,nil
}

func GetApplyrelations(receiverId string,status int)([]UserRelationShip,error)  {
	var relations []UserRelationShip
	rows , err := Database.Query("SELECT user_id, f_user_id from im_user_relation WHERE  receiver_id = ? AND status = ?",receiverId,status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var re UserRelationShip
		rows.Scan(&re.UserId,&re.FUserId)
		relations = append(relations, re)
	}
	return relations,nil
}


func GetApplyFriendList(userId string)([]SimpleUser,error)  {
	relations ,err := GetApplyrelations(userId,0)
	if err != nil {
		return nil, err
	}
	var users []SimpleUser
	for _, re := range relations {
		tempId := re.UserId
		if strings.EqualFold(re.UserId,userId) {
			tempId = re.FUserId
		}
		su,_ := GetUserWithLocationByUserId(tempId)
		if err == nil {
			users = append(users, *su)
		}
	}
	return users,nil
}