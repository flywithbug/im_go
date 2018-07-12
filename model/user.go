package model

import (
	_ "database/sql"
	"github.com/pborman/uuid"
	"encoding/json"
	"time"
	log "github.com/flywithbug/log4go"
)

type SimpleUser struct {
	Uid       	int64     	`json:"u_id"`        //id 就是Uid
	appId    	int64	   	`json:"app_id"`

	UserId 	 	string	   	`json:"user_id"`   //uuid生成
	Nick     	string    	`json:"nick"`      //昵称
	Status   	int    		`json:"status"`    //状态 0离线,1在线
	Sign     	string    	`json:"sign"`      //个性签名
	Avatar   	string    	`json:"avatar"`    //头像
	Forbidden 	int32		`json:"forbidden"`


}

/*
用户对象
*/
type User struct {
	SimpleUser
	createAt 	time.Time 	`json:"create_at"` //注册日期
	updateAt 	time.Time 	`json:"update_at"` //更新日期
	Token    	string    	`json:"token"`
	userAgent   string		`json:"user_agent"` //用户登录时的ua
}

func (use *User)GetAppId()int64  {
	return use.appId
}


/*
 转JSON数据
*/
func (this *User) Encode() []byte {
	s, _ := json.Marshal(*this)
	return s
}

/*
 解析JSON
*/
func (this *User) Decode(data []byte) error {
	err := json.Unmarshal(data, this)
	return err
}

/*
 检查账号是否存在
*/
func CheckAccount(account string) (int, error) {
	var num int
	rows, err := Database.Query("select count(*) from im_user where account=? ", account)
	if err != nil {
		log.Error(err.Error())
		return -1, &DatabaseError{"根据账号查询用户错误"}
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&num)
	}
	return num, nil
}

/*
 根据ID获取用户
*/
func GetUserByUserId(userId string) (*User, error) {
	var user User
	row := Database.QueryRow("select id, app_id, user_id, nick, status, sign, avatar, create_at, update_at from im_user where user_id = ?", userId)
	err := row.Scan(&user.Uid,&user.appId,&user.UserId, &user.Nick, &user.Status, &user.Sign, &user.Avatar, &user.createAt, &user.updateAt)
	if err != nil {
		log.Error(err.Error(),userId)
		return nil, &DatabaseError{"根据ID查询用户-将结果映射至对象错误"}
	}
	return &user, err
}
/*
 根据ID获取用户
*/
func GetUserByUId(uId string) (*User, error) {
	var user User
	row := Database.QueryRow("select id, app_id, user_id, nick, status, sign, avatar, create_at, update_at from im_user where id = ?", uId)
	err := row.Scan(&user.Uid,&user.appId,&user.UserId, &user.Nick, &user.Status, &user.Sign, &user.Avatar, &user.createAt, &user.updateAt)
	if err != nil {
		log.Error(err.Error()+uId)
		return nil, &DatabaseError{"根据ID查询用户-将结果映射至对象错误"}
	}
	return &user, err
}

/*
 根据token获取用户
*/
func GetUserByToken(token string) (*User, error) {
	var user User
	row := Database.QueryRow("select u.id, u.app_id,u.user_id,u.nick, u.status, u.sign, u.avatar, u.create_at, u.update_at from  im_user u left join im_login l on u.user_id=l.user_id where l.token=? AND l.status=1", token)
	err := row.Scan(&user.Uid,&user.appId,&user.UserId, &user.Nick, &user.Status, &user.Sign, &user.Avatar, &user.createAt, &user.updateAt)
	if err != nil {
		log.Error(err.Error()+token)
		return nil, &DatabaseError{"根据Token查询用户-将结果映射至对象错误"}
	}
	return &user, nil
}


/*
 登录账号
*/
func LoginUser(account string, password string) (*User, error) {
	var user User
	rows, err := Database.Query("select app_id, id , user_id, nick, status, sign, avatar, create_at, update_at,forbidden from im_user where account=? and password=? ", account, password)
	if err != nil {
		log.Error(err.Error()+account)
		return nil, &DatabaseError{"根据账号及密码查询用户错误"}
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.appId,&user.Uid ,&user.UserId, &user.Nick, &user.Status, &user.Sign, &user.Avatar, &user.createAt, &user.updateAt,&user.Forbidden)
		if err != nil {
			log.Error(err.Error())
			return nil, &DatabaseError{"根据账号及密码查询结果映射至对象错误"}
		}
	}
	return &user, nil
}

/*
 保存用户
*/
func SaveUser(appId int64,account string, password,origin_password string, nick string, avatar string) (*string, error) {
	insStmt, err := Database.Prepare("insert into im_user (user_id,app_id, account, password,origin_password, nick, avatar, create_at, update_at) VALUES (?,?, ?, ?, ?, ?, ?, ?,?)")
	if err != nil {
		log.Error(err.Error())
		return nil, &DatabaseError{"保存用户数据库处理错误"}
	}
	defer insStmt.Close()
	now := time.Now().Format("2006-01-02 15:04:05")
	uid := uuid.New()
	_, err = insStmt.Exec(uid, appId , account, password,origin_password, nick, avatar, now, now)
	if err != nil {
		log.Error(err.Error())
		return nil, &DatabaseError{"保存用户记录错误"}
	}
	return &uid, nil
}


/*
 根据条件查询获取好友列表
*/
func QueryUser(nick string) ([]SimpleUser, error) {
	var users []SimpleUser

	rows, err := Database.Query("SELECT id,user_id,nick,status,sign,avatar,forbidden FROM im_user WHERE nick LIKE ?",nick)
	if err != nil {
		return users, &DatabaseError{"根据查询用户错误"+err.Error()}
	}
	for rows.Next() {
		var user SimpleUser
		err =rows.Scan(&user.Uid,&user.UserId, &user.Nick, &user.Status, &user.Sign, &user.Avatar,&user.Forbidden)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateUserAvatar(avatar , userId string)error  {
	updateStmt,err := Database.Prepare("UPDATE im_user SET `avatar` = ? WHERE user_id=?")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res ,err := updateStmt.Exec(avatar,userId)
	if err != nil {
		return &DatabaseError{"服务出错"+err.Error()}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return  &DatabaseError{"未查询到该用户"}
	}
	return nil

}

func UpdateUserNickName(nick , userId string)error  {
	updateStmt,err := Database.Prepare("UPDATE im_user SET `nick` = ? WHERE user_id=?")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res ,err := updateStmt.Exec(nick,userId)
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务出错"+err.Error()}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return  &DatabaseError{"未查询到该用户"}
	}
	return nil

}

func UpdateuserPassWorld(old_password,password,origin_password,userId string)error  {
	updateStmt,err := Database.Prepare("UPDATE im_user SET `password` = ?,`origin_password`= ?  WHERE user_id=? AND origin_password = ? ")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	defer updateStmt.Close()
	res ,err := updateStmt.Exec(password,origin_password,userId,old_password)
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return  &DatabaseError{"未查询到该用户"}
	}
	return nil


}






