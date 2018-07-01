package model



type PushModel struct {
	Title 				string 	`json:"title"`
	SubTitle 			string	`json:"sub_title"`
	BadgeNumber 		int		`json:"badge_number"`
	AppId       		int		`json:"app_id"` //用于索引证书
	DeviceToken 		string	`json:"device_token"`
	EnvironmentType		int		`json:"environment_type"` //默认0位production环境
	Sound 				string	`json:"sound"`
	Body            	string	`json:"body"` //推送时显示的内容
	ContentAvailable	bool	`json:"content_available"` //1 表示正常推送，0表示静默推送

}

type MessageBoddy struct {
	Content   		string		`json:"content"`
	Type 			int			`json:"type"`
}