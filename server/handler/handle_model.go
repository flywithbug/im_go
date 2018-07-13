package handler

type loginoutModel struct {
	LocationModel
	//register
	AppId      int64	`json:"app_id"`
	Nick		string	`json:"nick"`
	Avatar      string	`json:"avatar"`

	//login
	Account 	string 	`json:"account"`
	Password 	string	`json:"password"`
	Signature 	string	`json:"signature"`     //account 和 password 拼接之后的加密字符串 用于非对称匹配

	//logout
	Token      string	`json:"token"`
	UserId     string    `json:"user_id"`    // 用户ID

	DeviceId  string	`json:"device_id"`

	UserAgent string    `json:"user_agent"`


	NewPassword string  `json:"new_password"`
	OldPassword	string	`json:"old_password"`


}

type LocationModel struct {
	Latitude    string		`json:"latitude"`   //维度
	Longitude   string		`json:"longitude"`   //经度
	LTimeStamp  string		`json:"l_time_stamp"`
	LType 		int			`json:"l_type"`
	PIdentifier	string		`json:"p_identifier"`
}

type PhotoLocationsModel struct {
	List   []LocationModel   `json:"list"`
}



//type relationShipModel struct {
//	UId 		int			`json:"u_id"`
//	FriendId	int			`json:"friend_id"`
//	Status 		int			`json:"status"`
//	RelationId 	string		`json:"relation_id"`
//	Remark		string		`json:"remark"`
//	Method 		string		`json:"method"`  //add delete
//}


