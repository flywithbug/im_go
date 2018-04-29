package handle



type loginoutModel struct {
	//login
	Account 	string 	`json:"account"`
	Password 	string	`json:"password"`
	Nick		string	`json:"nick"`
	Avatar      string	`json:"avatar"`


	//logout
	Token      string	`json:"token"`
	UserId    string    `json:"user_id"`    // 用户ID
}