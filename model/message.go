package model

import (
	"encoding/json"
	_ "database/sql"
	log "github.com/flywithbug/log4go"
	"github.com/pborman/uuid"
	"time"
)

const (
	IMMessageTypeText	= 1
	IMMessageTypePhoto	= 2
	IMMessageTypeAudio	= 3
	IMMessageTypeVideo	= 4
)

//Store in mysql
type IMMessage struct {
	Id				int				`json:"id"`
	MsgId			string			`json:"msg_id"`
	
	SMsgId			string			`json:"s_msg_id"`
	TimeStamp 		int				`json:"time_stamp"`
	Sender			int				`json:"sender"`
	Receiver		int				`json:"receiver"`
	Status 			int				`json:"status"`
	Type 			int				`json:"type"`
	Content 		string			`json:"content"`
	Font 			string			`json:"font"`
}

/*
 转JSON数据
 */
func (msg *IMMessage) Encode() []byte {
	s, _ := json.Marshal(msg)
	return s
}

/*
 解析JSON数据
 */
func (msg *IMMessage) Decode(data []byte) error {
	err := json.Unmarshal(data, msg)
	return err
}


//msgId 客户端生成的uuid  字符串长度36
func SaveIMMessage(sender ,receiver,msgType,status ,font , content, sMsgId string)(*string, error)  {
	insStmt ,err := Database.Prepare("INSERT into im_message (msg_id,sender,receiver,mtype,status,font,content,s_msg_id)VALUES (?,?,?,?,?,?,?,?) ")
	defer insStmt.Close()
	if err != nil {
		return nil,&DatabaseError{"消息服务出错"}
	}
	if len(sMsgId) == 0 {
		return nil,&DatabaseError{"消息id 为空"}
	}
	msgId := uuid.New()
	_,err = insStmt.Exec(msgId,sender,receiver,msgType,status,font,content,sMsgId)
	if err != nil{
		return nil,&DatabaseError{"保存消息错误"}
	}
	return &msgId,nil
}

func UpdateMessageACK(msgId string, status int)error  {
	updatStmt,err := Database.Prepare("UPDATE im_message SET `status` = ? `update_at`= ? WHERE msg_id = ?")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	res,err := updatStmt.Exec(status,time.Now().Unix(),msgId)
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return &DatabaseError{"未有记录被修改"}
	}
	return nil
}




