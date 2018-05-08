package model

import (
	"encoding/json"
	_ "database/sql"
	log "github.com/flywithbug/log4go"
	"time"
	"fmt"
)


//Store in mysql
type IMMessage struct {
	Id				int32				`json:"id"`			//msgId
	Sender			int32				`json:"sender"`
	Receiver		int32				`json:"receiver"`
	TimeStamp 		int32				`json:"time_stamp"`
	Status 			int32				`json:"status"`
	UpdateAt		int32				`json:"update_at"`
	Content 		[]byte				`json:"content"`   //客户端自行解析内容
}

func (msg *IMMessage) Description() string {
	return fmt.Sprintf("sender:%d,receiver:%d,timestamp:%d,msgId:%d,body:%s",msg.Sender,msg.Receiver,
		msg.TimeStamp,msg.Id,msg.Content)
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



//msgId 客户端生成的uuid  字符串长度36 返回数据库Id
func SaveIMMessage(sender, receiver, timestamp int32, body []byte)(int32, error)  {
	insStmt ,err := Database.Prepare("INSERT into im_message (sender,receiver,content,time_stamp)VALUES (?, ?, ?, ?) ")
	if err != nil {
		log.Error(err.Error())
		return -1,&DatabaseError{"消息服务出错"}
	}
	defer insStmt.Close()
	if sender == 0 || receiver == 0{
		return -1,&DatabaseError{"error parameter"}
	}
	if timestamp == 0 {
		timestamp = int32(time.Now().Unix())
	}
	res,err := insStmt.Exec(sender,receiver,body,timestamp)
	if err != nil{
		log.Error(err.Error())
		return -1,&DatabaseError{"保存消息错误"}
	}
	id,err := res.LastInsertId()
	if err != nil{
		log.Error(err.Error())
		return -1,&DatabaseError{"后去消息Id错误"}
	}
	return int32(id),nil
}

//发送状态回执
func UpdateMessageACK(msgId int32, status int)error  {
	updatStmt,err := Database.Prepare("UPDATE im_message SET `status` = ? ,`update_at`= ? WHERE id = ?")
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	defer updatStmt.Close()
	res,err := updatStmt.Exec(status,time.Now().Unix(),msgId)
	if err != nil {
		log.Error(err.Error())
		return  &DatabaseError{"服务出错"}
	}
	num, err := res.RowsAffected()
	if err != nil || num <= 0{
		return &DatabaseError{"未有记录被修改"}
	}
	return nil
}

//查找接收人消息
func FindeMessagesReceiver(receiver int32,status int)([]IMMessage,error)  {
	var messages []IMMessage
	rows ,err := Database.Query("SELECT id,sender,receiver,content,time_stamp,status,update_at FROM im_message WHERE receiver = ? AND status = ?",receiver,status)
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

//查找发送人消息
func FindeMessagesSender(sender int32,status int)([]IMMessage,error)  {
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




