package handler

type loginoutModel struct {

	//register
	AppId      int64	`json:"app_id"`
	Nick		string	`json:"nick"`
	Avatar      string	`json:"avatar"`

	//login
	Account 	string 	`json:"account"`
	Password 	string	`json:"password"`
	Key 		string	`json:"key"`     //account 和 password 拼接之后的加密字符串 用于非对称匹配

	//logout
	Token      string	`json:"token"`
	UserId     string    `json:"user_id"`    // 用户ID


	UserAgent string    `json:"user_agent"`
}

//type relationShipModel struct {
//	UId 		int			`json:"u_id"`
//	FriendId	int			`json:"friend_id"`
//	Status 		int			`json:"status"`
//	RelationId 	string		`json:"relation_id"`
//	Remark		string		`json:"remark"`
//	Method 		string		`json:"method"`  //add delete
//}


