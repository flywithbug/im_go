package model

import (
	"encoding/json"
	_ "database/sql"
	log "github.com/flywithbug/log4go"
	"time"
	"fmt"
)

const  (
	IMMessageReceiver_Type_offline  		=  1  //未发送成功，对方不在线
	IMMessageReceiver_Type_ClientACK  		=  2  //客户端收到，
	IMMessageReceiver_Type_Read  			=  3
	IMMessageReceiver_Type_ReCall  			=  4  //撤回
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
	Operation 		int32				`json:"operation"`  //消息类型，用于存储offline msg
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



////msgId 客户端生成的uuid  字符串长度36 返回数据库Id 转存到已发送的表中 删除offline表中的数据
//func SaveIMMessage(msg *IMMessage)(int32, error)  {
//	insStmt ,err := Database.Prepare("INSERT into im_message (id,sender,receiver,content,time_stamp)VALUES (?,?, ?, ?, ?) ")
//	if err != nil {
//		log.Error(err.Error())
//		return -1,&DatabaseError{"消息服务出错"}
//	}
//	defer insStmt.Close()
//	if msg.Sender == 0 || msg.Receiver == 0{
//		return -1,&DatabaseError{"error parameter"}
//	}
//	if msg.TimeStamp == 0 {
//		msg.TimeStamp = int32(time.Now().Unix())
//	}
//	_,err = insStmt.Exec(msg.Id,msg.Sender,msg.Receiver,msg.Content,msg.TimeStamp)
//	if err != nil{
//		log.Error(err.Error())
//		return -1,&DatabaseError{"保存消息错误"}
//	}
//	return msg.Id,nil
//}

func SaveIMMessageFromOfflineMessage(msgId int32)error  {
	insStmt ,err := Database.Prepare("INSERT into im_message (id,sender,receiver,content,time_stamp) SELECT id,sender,receiver,content,time_stamp from im_offline_message where id = ? ")
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"消息服务出错"}
	}
	defer insStmt.Close()
	insStmt.Exec(msgId)
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务出错"}
	}
	DeleteOfflineMessageMsgId(msgId)
	return UpdateMessageACK(msgId,1)
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



//msgId 客户端生成的uuid  字符串长度36 返回数据库Id
func SaveOfflineIMMessage(sender, receiver, timestamp,operation int32, body []byte)(int32, error)  {
	insStmt ,err := Database.Prepare("INSERT into im_offline_message (sender,receiver,content,time_stamp,operation)VALUES (?, ?, ?, ?,?) ")
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
	res,err := insStmt.Exec(sender,receiver,body,timestamp,operation)
	if err != nil{
		log.Error(err.Error())
		return -1,&DatabaseError{"保存消息错误"}
	}
	id,err := res.LastInsertId()
	if err != nil{
		log.Error(err.Error())
		return -1,&DatabaseError{"获取消息Id错误"}
	}
	return int32(id),nil
}


//查找接收人消息
func FindeOfflineMessagesReceiver(receiver int32,status int)([]IMMessage,error)  {
	var messages []IMMessage
	rows ,err := Database.Query("SELECT id,sender,receiver,content,time_stamp,status,update_at,operation FROM im_offline_message WHERE receiver = ? AND status = ?",receiver,status)
	defer rows.Close()
	if err != nil {
		log.Error(err.Error())
		return messages, &DatabaseError{"服务出错"}
	}
	for rows.Next(){
		var msg IMMessage
		rows.Scan(&msg.Id,&msg.Sender,&msg.Receiver,&msg.Content,&msg.TimeStamp,&msg.Status,&msg.UpdateAt,&msg.Operation)
		messages = append(messages, msg)
	}
	return messages,nil
}




func DeleteOfflineMessageMsgId(msgId int32)error  {
	delStmt,err := Database.Prepare("DELETE from  im_offline_message where  id = ?")
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务错误"}
	}
	defer  delStmt.Close()
	_,err = delStmt.Exec(msgId)
	if err != nil {
		log.Error(err.Error())
		return &DatabaseError{"服务错误"}
	}
	return nil
}
