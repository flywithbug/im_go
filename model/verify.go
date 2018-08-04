package model

import (
		_ "database/sql"
	"github.com/flywithbug/log4go"
	uuid "github.com/pborman/uuid"
	"math/rand"
	"time"
	"fmt"
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
func GeneryVerifyData(userId,account string,vld ,VType  int) (string,string,error) {
	insStmt, errStmt := Database.Prepare("insert into im_verify_code (uuid,verify,vld,v_type,user_id,account) VALUES (?, ?, ?, ?, ?,?)")
	if errStmt != nil {
		log4go.Info(errStmt.Error())
		return "","",&DatabaseError{"服务错误"}
	}
	defer insStmt.Close()
	uuId := uuid.NewUUID().String()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%04v", rnd.Int31n(10000))
	_, err := insStmt.Exec(uuId,vCode,vld,VType,userId,account)
	if err != nil {
		log4go.Info(err.Error())
		return "","",&DatabaseError{"服务错误"}
	}
	return uuId,vCode,nil
}

func CheckVerify(uuid string ,vType string) (userId,uuId string, err error)  {
	row := Database.QueryRow("select user_id ,uuid from im_verify_code where uuid=? and v_type=? and status= 0 ", uuid, vType)
	err = row.Scan(&userId,&uuid)
	if err != nil {
		log4go.Info(err.Error()+userId)
		return userId,uuId, &DatabaseError{"未查询到该验证信息"}
	}
	updateVerifyCodeStatus(uuId,1)
	return userId,uuId, nil
}

func CheckVerifyByAccount(account ,verify string,VType int) (useId,uuId string, err error)  {
	var  vld int
	row := Database.QueryRow("select user_id ,vld,uuid from im_verify_code where account=? and v_type=? and verify = ? and status= 0", account, VType,verify)
	err = row.Scan(&useId,&vld,&uuId)
	//if vld <  int(time.Now().Unix()){
	//	return "",&DatabaseError{"验证码超时未使用"}
	//}
	if err != nil {
		log4go.Info(err.Error()+account)
		return useId,uuId, &DatabaseError{"验证码不正确"}
	}
	updateVerifyCodeStatus(uuId,1)
	return useId,uuId,nil
}

func updateVerifyCodeStatus(uuId string, status int)  {
	updateStmt,err := Database.Prepare("UPDATE im_verify_code SET `status` = ?  WHERE uuid=?")
	if err != nil {
		log4go.Info(err.Error())
		return
	}
	defer updateStmt.Close()
	updateStmt.Exec(status,uuId)
}



