
//项目说明：存储部分使用的mysql, 暂未增加redis支持
##流程图
![流程图](http://7xsdes.com1.z0.glb.clouddn.com/1526641708662.jpg)

## 安装说明

//编译为centos 执行文件
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build main.go


系统需要安装Go和MySQL。

打开配置文件 config.json，修改相关配置。


创建数据库im_go，再导入install.sql

	$ mysql -u username -p -D im_go < install.sql






//api  
host localhost:8080

|method|	path|
|-------|------|
|Get| /api/system | 系统状态|
|Post| /api/register| 注册|
|Post|	login。|
|Post。|	/query。|



```
system  
{
    "status": 0,
    "msg": "",
    "data": {
        "cpu.info": [
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0
        ],
        "cpu.num": 8,
        "mem.free": "2460MB",
        "mem.total": "16384MB",
        "mem.used_percent": "67.31%"
    },
    "refer": ""
}

register
{
	"account":"ori2",
	"password":"ori",
	"avatar":"http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg",
	"nick":"ori"
}

{
    "status": 0,
    "msg": "注册成功",
    "data": null,
    "refer": ""
}

login
{
	"account":"ori2",
	"password":"ori"
}


{
    "status": 200,
    "msg": "",
    "data": {
        "user": {
            "id": 10003,
            "user_id": "fe9609b4-d22e-4672-b27c-1c13a3849f37",
            "nick": "ori",
            "status": "0",
            "sign": "",
            "avatar": "http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg",
            "token": "3680e5ef-d5e8-4efe-82c0-9ad03cf16025",
            "forbidden": 0
        }
    },
    "refer": "LOGIN_RETURN"
}

/query
{
	"nick":"ori"
}

{
    "status": 200,
    "msg": "",
    "data": {
        "users": [
            {
                "id": 10001,
                "user_id": "d5f75fbc-4f64-4f78-b320-2ca770847800",
                "nick": "ori",
                "status": 0,
                "sign": "",
                "avatar": "http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg",
                "forbidden": 0
            },
            {
                "id": 10002,
                "user_id": "80d07eb7-09d5-4332-aa4d-01990a291dfd",
                "nick": "ori",
                "status": 0,
                "sign": "",
                "avatar": "http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg",
                "forbidden": 0
            },
            {
                "id": 10003,
                "user_id": "fe9609b4-d22e-4672-b27c-1c13a3849f37",
                "nick": "ori",
                "status": 0,
                "sign": "",
                "avatar": "http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg",
                "forbidden": 0
            }
        ]
    },
    "refer": ""
}

relation
{
	"u_id":10001,
	"friend_id":10002,
	"method":"add"  
}

{
	"relation_id":"0b83e3965f12affa4371beaa267c3071",
	"remark":"232222",
	"method":"remark"
}

{
    "status": 200,
    "msg": "",
    "data": {
        "msg": "好友请求发送成功"
    },
    "refer": ""
}








```
