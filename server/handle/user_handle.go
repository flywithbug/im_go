package handle

import (
	"net/http"
	"im_go/model"
	"strings"
	"github.com/pborman/uuid"
	"im_go/common"
	"io/ioutil"
	"encoding/json"
	log "github.com/flywithbug/log4go"
)

// 注册请求
/*
	Para:appId,account,password,Nick，头像地址。
*/
func handleRegister(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body,_ := ioutil.ReadAll(req.Body)
		login := loginoutModel{}
		if err := json.Unmarshal(body,&login);err == nil {
			password := common.Md5(login.Password)
			register(resp, login.AppId,login.Account, password, login.Nick, login.Avatar)
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
		log.Debug("req.Body",string(body))
		m := loginoutModel{}
		if err := json.Unmarshal(body,&m);err == nil{
			password := common.Md5(m.Password)
			login(resp,m.Account,password,ip)
		}else {
			log.Error(err.Error())
			resp.Write(model.NewIMResponseSimple(401, "Bad Request: "+req.Method, "").Encode())
		}
	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}

func handleLogout(resp http.ResponseWriter,req*http.Request)  {
	// POST请求 退出登录
	if req.Method == "POST" {
		body,_ := ioutil.ReadAll(req.Body)
		m := loginoutModel{}
		if err := json.Unmarshal(body,&m);err == nil{
			logout(resp,m.Token)
		}else {
			log.Error(err.Error())
		}
	}else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}


/**
查询请求处理方法
*/
func handleQuery(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body,err := ioutil.ReadAll(req.Body)
		if err != nil{
			log.Error(err.Error())
			resp.Write(model.NewIMResponseSimple(401, "bad Request", "").Encode())
			return
		}
		m := loginoutModel{}
		if err := json.Unmarshal(body,&m);err == nil{
			users, _ := model.QueryUser(m.Nick)
			resp.Write(model.NewIMResponseData(common.SaveMapData("users", users), "").Encode())
		}else {
			resp.Write(model.NewIMResponseSimple(401, "bad Request", "").Encode())
		}
	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}


/*
	三个参数，
	u_id 好友请求发送者
	friend_id 接收人
	method  "add",delete,"remark" 添加备注
*/
func handleAddRelation(resp http.ResponseWriter, req *http.Request)  {
	if req.Method == "POST" {
		body,err := ioutil.ReadAll(req.Body)
		if err != nil{
			resp.Write(model.NewIMResponseSimple(401, "bad Request", "").Encode())
			return
		}
		m := relationShipModel{}
		if err := json.Unmarshal(body,&m);err == nil{
			//添加好友
			if strings.EqualFold(m.Method,"add")  {
				_, err = model.AddUserRelation(m.UId,m.FriendId)
				if err != nil {
					resp.Write(model.NewIMResponseSimple(500, "Server error", "").Encode())
				}else {
					resp.Write(model.NewIMResponseData(common.SaveMapData("msg", "success"), "").Encode())
				}
			}else if strings.EqualFold(m.Method,"delete") { //删除好友
				err = model.DelRelationShip(m.RelationId)
				if err != nil {
					resp.Write(model.NewIMResponseSimple(500, "Server error", "").Encode())
				}else {
					resp.Write(model.NewIMResponseData(common.SaveMapData("msg", "success"), "").Encode())
				}
			}else if strings.EqualFold(m.Method,"remark") { //删除好友
				err = model.UpdateRelationRemark(m.RelationId,m.Remark)
				if err != nil {
					resp.Write(model.NewIMResponseSimple(500, "Server error", "").Encode())
				}else {
					resp.Write(model.NewIMResponseData(common.SaveMapData("msg", "success"), "").Encode())
				}
			}else {
				resp.Write(model.NewIMResponseSimple(401, "bad Request", "").Encode())
			}
		}else {
			log.Error(err.Error())
			resp.Write(model.NewIMResponseSimple(401, "bad Request", "").Encode())
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
			if !strings.EqualFold(user.UserId, "") {
				token := uuid.New()
				if err := model.SaveLogin(user.GetAppId(),user.Uid,user.UserId, token, ip,user.Forbidden); err != nil {
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
func register(resp http.ResponseWriter,appId int64, account string, password string, nick string, avatar string) {
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
			if appId == 0{
				appId = 10
			}
			_, err := model.SaveUser(appId,account, password, nick, avatar)
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