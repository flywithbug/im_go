package handle

import (
	"net/http"
	"im_go/model"
	"strings"
	"github.com/pborman/uuid"
	"im_go/common"
	"io/ioutil"
	"encoding/json"
)



// 注册请求
func handleRegister(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body,_ := ioutil.ReadAll(req.Body)
		login := loginoutModel{}
		if err := json.Unmarshal(body,&login);err == nil {
			password := common.Md5(login.Password)
			register(resp, login.Account, password, login.Nick, login.Avatar)
		}else {
			resp.Write(model.NewIMResponseSimple(401, "Bad Request: "+req.Method, "").Encode())
		}
	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}

/**
登录请求处理方法
*/
func handleLogin(resp http.ResponseWriter, req *http.Request) {
	// POST登录请求
	if req.Method == "POST" {
		ip := common.GetIp(req)
		body,_ := ioutil.ReadAll(req.Body)
		m := loginoutModel{}
		if err := json.Unmarshal(body,&m);err == nil{
			password := common.Md5(m.Password)
			login(resp,m.Account,password,ip)
		}else {
			resp.Write(model.NewIMResponseSimple(401, "Bad Request: "+req.Method, "").Encode())
		}
	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}

func handleLogout(resp http.ResponseWriter,req*http.Request)  {
	// POST登录请求
	if req.Method == "POST" {
		body,_ := ioutil.ReadAll(req.Body)
		m := loginoutModel{}
		if err := json.Unmarshal(body,&m);err == nil{
			logout(resp,m.Token)
		}
	}else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}




// 登录主方法
func login(resp http.ResponseWriter, account string, password string, ip string) {
	if account == "" {
		resp.Write(model.NewIMResponseSimple(401, "账号不能为空", "").Encode())
	} else if password == "" {
		resp.Write(model.NewIMResponseSimple(401, "密码不能为空", "").Encode())
	} else {
		num, err := model.CheckAccount(account)
		if err != nil {
			resp.Write(model.NewIMResponseSimple(500, err.Error(), "").Encode())
			return
		}
		if num > 0 {
			user, err := model.LoginUser(account, password)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(500, err.Error(), "").Encode())
				return
			}
			if !strings.EqualFold(user.Id, "") {
				token := uuid.New()
				if _, err := model.SaveLogin(user.Id, token, ip); err != nil {
					resp.Write(model.NewIMResponseSimple(500, err.Error(), "").Encode())
				} else {
					user.Token = token
					resp.Write(model.NewIMResponseData(common.SaveMapData("user", user), "LOGIN_RETURN").Encode())
				}
			} else {
				resp.Write(model.NewIMResponseSimple(401, "密码错误", "").Encode())
			}
		} else {
			resp.Write(model.NewIMResponseSimple(401, "账号不存在", "").Encode())
		}
	}
}

func logout(resp http.ResponseWriter,token string)  {
	if token == "" {
		resp.Write(model.NewIMResponseSimple(401, "token不能为空", "").Encode())
	}else {
		num ,err := model.Logout(token)
		if num <= 0 || err != nil{
			errStr := err.Error()
			resp.Write(model.NewIMResponseSimple(500, errStr, "").Encode())
		}else {
			resp.Write(model.NewIMResponseData(common.SaveMapData("msg","success"),"").Encode())
		}
	}
}



/*
 用户注册
 101	账号不能为空
 102	密码不能为空
 103	用户名已存在
 104	昵称不能为空
 105	注册失败
*/
func register(resp http.ResponseWriter, account string, password string, nick string, avatar string) {
	if account == "" {
		resp.Write(model.NewIMResponseSimple(101, "账号不能为空", "").Encode())
	} else if password == "" {
		resp.Write(model.NewIMResponseSimple(102, "密码不能为空", "").Encode())
	} else if nick == "" {
		resp.Write(model.NewIMResponseSimple(103, "昵称不能为空", "").Encode())
	} else {
		num, err := model.CheckAccount(account)
		if err != nil {
			resp.Write(model.NewIMResponseSimple(103, err.Error(), "").Encode())
			return
		}
		if num > 0 {
			resp.Write(model.NewIMResponseSimple(104, "用户名已存在", "").Encode())
		} else {
			_, err := model.SaveUser(account, password, nick, avatar)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(104, err.Error(), "").Encode())
				return
			}
			num, _ := model.CheckAccount(account)
			if num > 0 {
				resp.Write(model.NewIMResponseSimple(0, "注册成功", "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(105, "注册失败", "").Encode())
			}
		}
	}
}
