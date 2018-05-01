package model

import (
	"github.com/pborman/uuid"
	"time"
	_ "database/sql"
	"fmt"
)

const (
	STATUS_LOGIN  = 	int8(1)
	STATUS_LOGOUT  = 	int8(2)
)

type Login struct {
	Id        	string    	`json:"id"`         // id
	UserId    	string    	`json:"user_id"`    // 用户ID
	Token     	string    	`json:"token"`      // 用户TOKEN
	LoginAt   	time.Time 	`json:"login_at"`   // 登录日期
	LoginIp   	string    	`json:"login_ip"`   // 登录IP
	Status 	  	int8		`json:"status"`		//status 1 已登录，0表示退出登录
	Forbidden 	int32		`json:"forbidden"`  //false 表示未禁言
	AppId    	int32		`json:"app_id"`
	UId      	int64		`json:"u_id"`
}


/*
 根据token和用户登录状态获取用户登录
 */
func GetLoginByToken(token string,status int8) (*Login, error) {
	var login Login
	row := Database.QueryRow("select u_id, id, user_id, token, login_at, login_ip, forbidden  from im_login where token=? AND status = ?", token,status)
	err := row.Scan(&login.UId,&login.Id, &login.UserId, &login.Token, &login.LoginAt, &login.LoginIp,&login.Forbidden)
	if err != nil {
		return nil, &DatabaseError{"根据Token获取用户登录错误"}
	}
	return &login, nil
}


/*
 保存登录状态
 */
func SaveLogin(appId int64,uId int64 ,userId string, token string, ip string,forbidden int32) (*string, error) {
	insStmt, errStmt := Database.Prepare("insert into im_login (app_id,u_id,id, user_id, token, login_at, login_ip, status,forbidden) VALUES (?, ?, ?, ?, ?, ?, ?, ?,?)")
	if errStmt != nil {
		return nil, &DatabaseError{"服务错误"}
	}
	defer insStmt.Close()
	id := uuid.New()

	status := 1
	_, err := insStmt.Exec(appId,uId,id, userId, token, time.Now().Format("2006-01-02 15:04:05"), ip,status,forbidden)
	if err != nil {
		return nil, &DatabaseError{"服务错误"}
	}
	return &id, nil
}



/*
	退出登录
*/
func Logout(token string)(int64,error) {
	updateStmt,err := Database.Prepare("UPDATE im_login SET `status` = ? WHERE token=? AND status = 1")
	defer updateStmt.Close()
	if err != nil {
		fmt.Println(err)
		return -1, &DatabaseError{"服务出错"}
	}
	res ,err := updateStmt.Exec(0,token)
	if err != nil {
		return -1, &DatabaseError{"服务出错"}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return -1, &DatabaseError{"token已失效"}
	}
	return num,nil
}