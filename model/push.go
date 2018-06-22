package model



type PushModel struct {
	Title 			string 	`json:"title"`
	SubTitle 		string	`json:"sub_title"`
	BadgeNumber 	int	`json:"badge_number"`
	AppId       	int		`json:"app_id"` //用于索引证书
	DeviceToken 	string	`json:"device_token"`
	EnvironmentType	int		`json:"environment_type"` //默认0位production环境
	Sound 			string	`json:"sound"`
	Body            string	`json:"body"` //推送时显示的内容
}

type MessageBoddy struct {
	Content   		string		`json:"content"`
	Type 			int			`json:"type"`
}