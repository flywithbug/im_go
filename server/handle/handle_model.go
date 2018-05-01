package handle



type loginoutModel struct {

	//register
	AppId      int64	`json:"app_id"`
	Nick		string	`json:"nick"`
	Avatar      string	`json:"avatar"`

	//login
	Account 	string 	`json:"account"`
	Password 	string	`json:"password"`


	//logout
	Token      string	`json:"token"`
	UserId    string    `json:"user_id"`    // 用户ID

}