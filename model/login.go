package model

import (
	"time"
	_ "database/sql"
	log "github.com/flywithbug/log4go"
)

const (
	STATUS_LOGIN  = 	1
	STATUS_LOGOUT  = 	2
)

type Login struct {
	Id        	int64    	`json:"id"`         // id
	UserId    	string    	`json:"user_id"`    // 用户ID
	Token     	string    	`json:"token"`      // 用户TOKEN
	LoginAt   	time.Time 	`json:"login_at"`   // 登录日期
	LoginIp   	string    	`json:"login_ip"`   // 登录IP
	Status 	  	int			`json:"status"`		//status 1 已登录，0表示退出登录
	Forbidden 	int32		`json:"forbidden"`  //false 表示未禁言
	AppId    	int64		`json:"app_id"`
	UId      	int32		`json:"u_id"`
}


/*
 根据token和用户登录状态获取用户登录
 */
func GetLoginByToken(token string) (*Login, error) {
	var login Login
	row := Database.QueryRow("SELECT id ,u_id,app_id, user_id, token, login_at, login_ip, forbidden,status  FROM im_login WHERE token=?", token)
	err := row.Scan(&login.Id,&login.UId,&login.AppId, &login.UserId, &login.Token, &login.LoginAt, &login.LoginIp,&login.Forbidden,&login.Status)
	if err != nil {
		log.Warn(err.Error())
		return nil, &DatabaseError{"根据Token获取用户登录错误"}
	}
	return &login, nil
}

//获取已登录的记录，非当前查询的token
//func GetLoginByUserId(userId,nToken string)(*User,error)  {
//
//
//}



/*
 保存登录状态
 */
func SaveLogin(appId int64,uId int64 ,userId string, token string, ip string,forbidden int32) error {
	insStmt, errStmt := Database.Prepare("insert into im_login (app_id,u_id, user_id, token, login_at, login_ip, status,forbidden) VALUES (?, ?, ?, ?, ?, ?, ?,?)")
	if errStmt != nil {
		return &DatabaseError{"服务错误"}
	}
	defer insStmt.Close()

	status := 1
	_, err := insStmt.Exec(appId,uId, userId, token, time.Now().Format("2006-01-02 15:04:05"), ip,status,forbidden)
	if err != nil {
		return &DatabaseError{"服务错误"}
	}
	return nil
}



/*
	退出登录  status 1 登录状态，0 是退出 -1 被其他登录用户踢出
*/
func Logout(token string)(int64,error) {
	updateStmt,err := Database.Prepare("UPDATE im_login SET `status` = ?,logout_at=? WHERE token=? AND status = 1")
	defer updateStmt.Close()
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"服务出错"}
	}

	res ,err := updateStmt.Exec(0,time.Now().Format("2006-01-02 15:04:05"),token)
	if err != nil {
		return -1, &DatabaseError{"服务出错"}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return -1, &DatabaseError{"token已失效"}
	}
	return num,nil
}

func LogoutOthers(token string,uid int32)(int64,error)  {
	updateStmt,err := Database.Prepare("UPDATE im_login SET `status` = ?,logout_at=? WHERE token<>? AND status = 1 AND u_id = ? ")
	defer updateStmt.Close()
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"服务出错"}
	}
	res ,err := updateStmt.Exec(0,time.Now().Format("2006-01-02 15:04:05"),token,uid)
	if err != nil {
		return -1, &DatabaseError{"服务出错"}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return -1, &DatabaseError{"token已失效"}
	}
	return num,nil
}