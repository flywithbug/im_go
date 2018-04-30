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


/**
查询请求处理方法
*/
func handleQuery(resp http.ResponseWriter, req *http.Request) {
	body,_ := ioutil.ReadAll(req.Body)
	m := loginoutModel{}
	if err := json.Unmarshal(body,&m);err == nil{
		users, _ := model.QueryUser("nick", "like", m.Nick)
		resp.Write(model.NewIMResponseData(common.SaveMapData("users", users), "").Encode())
	}else {
		resp.Write(model.NewIMResponseSimple(404, err.Error()+req.Method, "").Encode())

	}
	nick := req.FormValue("nick")
	users, err := model.QueryUser("nick", "like", nick)
	if err == nil {
		resp.Write(model.NewIMResponseData(common.SaveMapData("users", users), "").Encode())
	}
}


// 添加好友，获取好友列表
func handleUserCategoryAdd(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		//获取好友列表
		userId := req.FormValue("user_id")
		categories, err := model.GetCategoriesByUserId(userId)
		if err != nil {
			resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
			return
		}
		categories, err = model.GetBuddiesByCategories(categories)
		if err != nil {
			resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
			return
		}
		resp.Write(model.NewIMResponseData(common.SaveMapData("categories", categories), "").Encode())
	case "POST":
		// 添加好友列表
		userId := req.FormValue("user_id")
		name := req.FormValue("name")
		if userId == "" {
			resp.Write(model.NewIMResponseSimple(101, "用户ID不能为空", "").Encode())
		} else if name == "" {
			resp.Write(model.NewIMResponseSimple(102, "类别名称不能为空", "").Encode())
		} else {
			_, err := model.AddCategory(userId, name)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(103, err.Error(), "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(200, "添加分类成功", "").Encode())
			}
		}
	default:
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())

	}
}


// 删除好友分类
func handleUserCategoryDel(resp http.ResponseWriter, req *http.Request) {
	categoryId := req.FormValue("category_id")
	switch req.Method {
	case "GET":
		if categoryId == "" {
			resp.Write(model.NewIMResponseSimple(102, "类别ID不能为空", "").Encode())
		} else {
			num, err := model.DelCategoryById(categoryId)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			}
			if num > 0 {
				resp.Write(model.NewIMResponseSimple(200, "已删除好友分类", "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(103, "删除好友分类失败", "").Encode())
			}
		}
	case "POST":
		if categoryId == "" {
			resp.Write(model.NewIMResponseSimple(102, "类别ID不能为空", "").Encode())
		} else {
			num, err := model.DelCategoryById(categoryId)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			}
			if num > 0 {
				resp.Write(model.NewIMResponseSimple(200, "已删除好友分类", "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(103, "删除好友关系分类", "").Encode())
			}
		}
	default:
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}

// 编辑好友分类
func handleUserCategoryEdit(resp http.ResponseWriter, req *http.Request) {
	categoryId := req.FormValue("category_id")
	categoryName := req.FormValue("category_name")
	switch req.Method {
	case "GET":
		if categoryId == "" {
			resp.Write(model.NewIMResponseSimple(101, "类别ID不能为空", "").Encode())
		} else if categoryName == "" {
			resp.Write(model.NewIMResponseSimple(102, "类别名称不能为空", "").Encode())
		} else {
			num, err := model.EditCategoryById(categoryId, categoryName)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			}
			if num > 0 {
				resp.Write(model.NewIMResponseSimple(200, "修改用户好友类别成功", "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(103, "修改用户好友类别失败", "").Encode())
			}
		}
	case "POST":
		if categoryId == "" {
			resp.Write(model.NewIMResponseSimple(101, "类别ID不能为空", "").Encode())
		} else if categoryName == "" {
			resp.Write(model.NewIMResponseSimple(102, "类别名称不能为空", "").Encode())
		} else {
			num, err := model.EditCategoryById(categoryId, categoryName)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			}
			if num > 0 {
				resp.Write(model.NewIMResponseSimple(200, "修改用户好友类别成功", "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(103, "修改用户好友类别失败", "").Encode())
			}
		}
	default:
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())

	}
}
func handleUserCategoryQuery(resp http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	categories, err := model.GetCategoriesByUserId(id)
	if err != nil {
		resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
	} else {
		resp.Write(model.NewIMResponseData(common.SaveMapData("categories", categories), "").Encode())
	}
}

// 添加好友关系
func handleUserRelationAdd(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		receiver_category_id := req.FormValue("receiver_category_id")
		buddy_request_id := req.FormValue("buddy_request_id")
		buddyrequest, _ := model.GetBuddyRequestById(buddy_request_id)
		if buddyrequest != nil {
			receiver := buddyrequest.Receiver
			sender := buddyrequest.Sender
			sender_category_id := buddyrequest.SenderCategoryId
			//开启事务
			tx, _ := model.Database.Begin()
			//修改好友请求记录中接受人的好友分组ID
			_, err := model.UpdateBuddyRequestReceiverCategoryId(tx, buddy_request_id, receiver_category_id)
			//添加请求人好友关系数据
			_, err = model.AddFriendRelation(tx, receiver, sender_category_id)
			//添加接收人好友关系数据
			_, err = model.AddFriendRelation(tx, sender, receiver_category_id)
			//修改好友请求记录中状态
			_, err = model.UpdateBuddyRequestStatus(tx, buddy_request_id, "1")

			if err != nil {
				tx.Rollback()
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			} else {
				tx.Commit()
				//TODO 好友关系通知
				//判断请求者是不是在线 在线就把接受者推送给请求者
				//conn, _ := model.GetConnByUserId(sender)
				//if conn != nil { //在线
				//	user, _ := model.GetUserById(receiver)
				//	data := make(map[string]interface{})
				//	data["category_id"] = sender_category_id
				//	data["user"] = user
				//	ClientMaps[conn.Key].PutOut(common.NewIMResponseData(util.SetData("user", data), common.ADD_BUDDY))
				//}
				//conn, _ = model.GetConnByUserId(receiver)
				//if conn != nil {
				//	user, _ := model.GetUserById(sender)
				//	data := make(map[string]interface{})
				//	data["category_id"] = receiver_category_id
				//	data["user"] = user
				//	ClientMaps[conn.Key].PutOut(common.NewIMResponseData(util.SetData("user", data), common.ADD_BUDDY))
				//}
				//resp.Write(common.NewIMResponseSimple(0, "好友关系建立成功", "").Encode())
				return
			}

		} else {
			resp.Write(model.NewIMResponseSimple(104, "该好友请求不存在", "").Encode())
		}

	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}

// 删除好友关系
func handleUserRelationDel(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		userId := req.FormValue("user_id")
		categoryId := req.FormValue("category_id")
		if userId == "" {
			resp.Write(model.NewIMResponseSimple(101, "用户ID不能为空", "").Encode())
		} else if categoryId == "" {
			resp.Write(model.NewIMResponseSimple(102, "类别ID不能为空", "").Encode())
		} else {
			num, err := model.DelFriendRelation(userId, categoryId)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			}
			if num > 0 {
				resp.Write(model.NewIMResponseSimple(200, "已删除好友关系", "").Encode())
			} else {
				resp.Write(model.NewIMResponseSimple(103, "删除好友关系失败", "").Encode())
			}
		}
	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}
func handleUserRelationPush(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sender_category_id := req.FormValue("sender_category_id")
		sender := req.FormValue("sender")
		receiver := req.FormValue("receiver")
		if sender_category_id == "" {
			resp.Write(model.NewIMResponseSimple(101, "请选择分组", "").Encode())
		} else if sender == "" {
			resp.Write(model.NewIMResponseSimple(102, "请重新登录", "").Encode())
		} else {
			//TODO
			//判断接收人是不是在线 在线直接推送，不在线记录至请求表中
			//conn, _ := model.GetConnByUserId(receiver)
			//user, _ := model.GetUserById(sender)
			_, err := model.AddBuddyRequest(sender, sender_category_id, receiver)
			if err != nil {
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
			} else {
				//if conn != nil { //在线 直接推送 不在线 客户登录时候会激活请求通知
				//	data := make(map[string]interface{})
				//	data["id"] = user.Id
				//	data["nick"] = user.Nick
				//	data["status"] = user.Status
				//	data["sign"] = user.Sign
				//	data["avatar"] = user.Avatar
				//	data["buddyRequestId"] = buddyRequestId
				//	ClientMaps[conn.Key].PutOut(common.NewIMResponseData(util.SetData("user", data), common.PUSH_BUDDY_REQUEST))
				//	resp.Write(common.NewIMResponseSimple(0, "发送好友请求成功", "").Encode())
				//}
				resp.Write(model.NewIMResponseSimple(200, "发送好友请求成功", "").Encode())
				return
			}
		}
	} else {
		resp.Write(model.NewIMResponseSimple(404, "Not Found: "+req.Method, "").Encode())
	}
}

func handleUserRelationRefuse(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		buddy_request_id := req.FormValue("buddy_request_id")
		if buddy_request_id != "" {
			tx, _ := model.Database.Begin()
			//修改好友请求记录中状态
			_, err := model.UpdateBuddyRequestStatus(tx, buddy_request_id, "2")
			if err != nil {
				tx.Rollback()
				resp.Write(model.NewIMResponseSimple(100, err.Error(), "").Encode())
				return
			} else {
				tx.Commit()
				resp.Write(model.NewIMResponseSimple(200, "已经拒绝该好友请求成功", "").Encode())
				return
			}
		} else {
			resp.Write(model.NewIMResponseSimple(109, "该好友请求不合法", "").Encode())
		}

	} else {
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
