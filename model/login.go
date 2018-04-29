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
	insStmt, errStmt := Database.Prepare("insert into im_login (id, user_id, token, login_at, login_ip) VALUES (?, ?, ?, ?,?)")
	if errStmt != nil {
		return nil, &DatabaseError{"保存用户登录记录错误，数据库语句错误"}
	}
	defer insStmt.Close()
	id := uuid.New()
	_, err := insStmt.Exec(id, userId, token, time.Now().Format("2006-01-02 15:04:05"), ip)
	if err != nil {
		return nil, &DatabaseError{"保存用户登录记录错误"}
	}
	return &id, nil
}



/*
	退出登录
*/
func Logout(token string)(int64,error) {
	updateStmt,err := Database.Prepare("UPDATE im_login SET 'token' = ? WHERE token=?")
	defer updateStmt.Close()
	if err != nil {
		fmt.Println(err)
		return -1, &DatabaseError{"数据库处理失败"}
	}
	res ,err := updateStmt.Exec("",token)
	if err != nil {
		return -1, &DatabaseError{"更新token失败"}
	}
	num, err := res.RowsAffected()
	if err != nil {
		return -1, &DatabaseError{"读取token更新影响行数错误"}
	}
	return num,nil
}