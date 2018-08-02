package model

import (
		_ "database/sql"
	"github.com/flywithbug/log4go"
	uuid "github.com/pborman/uuid"
)

type Verify struct {
	UUID    	string   	`json:"uuid"`
	Verify  	string		`json:"verify"`
	Vld         int			`json:"vld"` //有效期
	VType       int			`json:"v_type"` //验证码类型
	Account     string      `json:"account"`
}


/*
 保存登录状态
 */
func GeneryVerifyData(verify,userId string,vld ,VType  int) (string,error) {
	insStmt, errStmt := Database.Prepare("insert into im_verify_code (uuid,verify,vld,v_type,user_id) VALUES (?, ?, ?, ?, ?)")
	if errStmt != nil {
		log4go.Info(errStmt.Error())
		return "",&DatabaseError{"服务错误"}
	}
	defer insStmt.Close()
	uuId := uuid.NewUUID().String()
	_, err := insStmt.Exec(uuId,verify,vld,VType,userId)
	if err != nil {
		log4go.Info(err.Error())
		return "",&DatabaseError{"服务错误"}
	}
	return uuId,nil
}

func CheckVerify(uuid string ,vType string) (user_id string, err error)  {
	row := Database.QueryRow("select user_id from im_verify_code where uuid=? and v_type=? ", uuid, vType)
	err = row.Scan(&user_id)
	if err != nil {
		log4go.Error(err.Error()+user_id)
		return user_id, &DatabaseError{"未查询到该验证信息"}
	}
	return user_id, nil
}



