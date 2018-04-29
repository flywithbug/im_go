package model

import (
	"github.com/pborman/uuid"
	"time"
	_ "database/sql"
	"fmt"
)

type Login struct {
	Id        string    `json:"id"`         // id
	UserId    string    `json:"user_id"`    // 用户ID
	Token     string    `json:"token"`      // 用户TOKEN
	LoginAt   time.Time `json:"login_at"`   // 登录日期
	LoginIp   string    `json:"login_ip"`   // 登录IP
}


/*
 根据token获取用户登录
 */
func GetLoginByToken(token string) (*Login, error) {
	var login Login
	row := Database.QueryRow("select id, user_id, token, login_at, login_ip from im_login where token=?", token)
	err := row.Scan(&login.Id, &login.UserId, &login.Token, &login.LoginAt, &login.LoginIp)
	if err != nil {
		return nil, &DatabaseError{"根据Token获取用户登录错误"}
	}
	return &login, nil
}



/*
 保存登录状态
 */
func SaveLogin(userId string, token string, ip string) (*string, error) {
	insStmt, errStmt := Database.Prepare("insert into im_login (id, user_id, token, login_at, login_ip, status) VALUES (?, ?, ?, ?, ?, ?)")
	if errStmt != nil {
		return nil, &DatabaseError{"服务错误"}
	}
	defer insStmt.Close()
	id := uuid.New()
	status := 1
	_, err := insStmt.Exec(id, userId, token, time.Now().Format("2006-01-02 15:04:05"), ip,status)
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
		return -1, &DatabaseError{"退出登录失败"}
	}
	num, err := res.RowsAffected()
	if err != nil{
		return -1, &DatabaseError{"服务出错"}
	}
	return num,nil
}