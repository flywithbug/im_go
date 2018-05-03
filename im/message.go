package im

import (
	"encoding/json"
)

type Message struct {
	Id				int				`json:"id"`
	Sender			int				`json:"sender"`
	Receiver		int				`json:"receiver"`
	Status 			int				`json:"status"`
	Type 			int				`json:"type"`
	TimeStamp 		int				`json:"time_stamp"`
	UpdateAt		int				`json:"update_at"`

	Content 		string			`json:"content"`
	Font 			string			`json:"font"`

	MsgId			string			`json:"msg_id"`
}


/*
 转JSON数据
 */
func (msg *Message) Encode() []byte {
	s, _ := json.Marshal(msg)
	return s
}

/*
 解析JSON数据
 */
func (msg *Message) Decode(data []byte) error {
	err := json.Unmarshal(data, msg)
	return err
}
