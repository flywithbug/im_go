package model

import (
	"encoding/json"
	_ "database/sql"
	log "github.com/flywithbug/log4go"
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
	Id				int				`json:"id"`			//msgId
	Sender			int				`json:"sender"`
	Receiver		int				`json:"receiver"`
	TimeStamp 		int				`json:"time_stamp"`
	Status 			int				`json:"status"`
	UpdateAt		int				`json:"update_at"`
	Content 		[]byte			`json:"content"`   //客户端自行解析内容
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



func MessageOperation(sender int,receiver int,content string)(msg *IMMessage,err error)  {



	return
}


//msgId 客户端生成的uuid  字符串长度36 返回数据库Id
func SaveIMMessage(sender ,receiver,msgType,status int, content, msgId string)(int64, error)  {
	insStmt ,err := Database.Prepare("INSERT into im_message (sender,receiver,mtype,status,content,msg_id,time_stamp)VALUES (?,?,?,?,?,?,?,?) ")
	defer insStmt.Close()
	if err != nil {
		return -1,&DatabaseError{"消息服务出错"}
	}
	if len(msgId) == 0 || sender == 0 || receiver == 0{
		return -1,&DatabaseError{"error parameter"}
	}
	now := time.Now().Unix()
	res,err := insStmt.Exec(sender,receiver,msgType,status,content,msgId,now)
	if err != nil{
		return -1,&DatabaseError{"保存消息错误"}
	}
	id,err := res.LastInsertId()
	if err != nil{
		return -1,&DatabaseError{"后去消息Id错误"}
	}
	return id,nil
}

//发送状态回执
func UpdateMessageACK(msgId string, status int)error  {
	updatStmt,err := Database.Prepare("UPDATE im_message SET `status` = ? `update_at`= ? WHERE msg_id = ?")
	defer updatStmt.Close()
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

func FindeMessages(sender string,status int)([]IMMessage,error)  {
	var messages []IMMessage
	rows ,err := Database.Query("SELECT id,sender,receiver,content,time_stamp,status,update_at FROM im_message WHERE sender = ? AND status = ?",sender,status)
	defer rows.Close()
	if err != nil {
		log.Error(err.Error())
		return messages, &DatabaseError{"服务出错"}
	}
	for rows.Next(){
		var msg IMMessage
		rows.Scan(&msg.Id,&msg.Sender,&msg.Receiver,&msg.Content,&msg.TimeStamp,&msg.Status,&msg.UpdateAt)
		messages = append(messages, msg)
	}
	return messages,nil
}





