package server




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

type relationShipModel struct {
	UId 		int			`json:"u_id"`
	FriendId	int			`json:"friend_id"`
	Status 		int			`json:"status"`
	RelationId 	string		`json:"relation_id"`
	Remark		string		`json:"remark"`
	Method 		string		`json:"method"`  //add delete
}


