package model

import (
	"encoding/json"
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
	Sender			int				`json:"sender"`
	Receiver		int				`json:"receiver"`
	SendAt 			time.Time		`json:"send_at"`
	Content 		string			`json:"content"`
	Type 			int				`json:"type"`
	Font 			string			`json:"font"`
	Status 			int				`json:"status"`
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

func SaveIMMessage(sender ,receiver,msgType int,send_at time.Time,font string)  {

}





