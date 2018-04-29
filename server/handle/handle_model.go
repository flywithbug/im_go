package handle



type loginModel struct {
	Account 	string 	`json:"account"`
	Password 	string	`json:"password"`
	Nick		string	`json:"nick"`
	Avatar      string	`json:"avatar"`
}