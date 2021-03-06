package handler

type loginoutModel struct {
	LocationModel
	//register
	AppId      int64	`json:"app_id"`
	Nick		string	`json:"nick"`
	Avatar      string	`json:"avatar"`
	Mail        string  `json:"mail"`



	//login
	Account 		string 	`json:"account"`
	Password 		string	`json:"password"`
	VerifyKey	    string	`json:"verify_key"` //验证码IdKey
	Verify	 		string	`json:"verify"`   //验证码


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


type UserRelationShip struct {
	UserId   	string			`json:"user_id"`
	Status 		int				`json:"status"`  //-2拉黑，-1 拒绝，0 申请，1 接受
}





