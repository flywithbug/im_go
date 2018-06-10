package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

var config *IMConfig

func Conf() *IMConfig {
	return config
}

/*
IM配置结构体
*/
type IMConfig struct {
	IMPort     		int      	`json:"im_port"`     	//IM服务长连接监听端口
	HttpPort   		int      	`json:"http_port"`   	//IM服务外部调研接口
	LogLevel        int			`json:"log_level"`    // 0:release  1:debug
	ServerPort 		string		`json:"server_port"`	//用户关系相关服务
	ProfBind		string	 	`json:"prof_bind"`		//机器监控
	DBConfig   		DBConfig 	`json:"db_config"`   	//数据库配置
	RouterPrefix 	[]string 	`json:"router_prefix"` //api前缀
	AuthFilterWhite []string 	`json:"auth_filter_white"` //api前缀
	AppConfig		AppConfig	`json:"app_config"`
}

/*
数据库配置结构体
*/
type DBConfig struct {
	Host         string `json:"host"`           //连接地址
	Username     string `json:"username"`       //用户名
	Password     string `json:"password"`       //用户密码
	Name         string `json:"name"`           //数据库名
	MaxIdleConns int    `json:"max_idle_conns"` //连接池最大空闲连接数
	MaxOpenConns int    `json:"max_open_conns"` //连接池最大连接数
}

type AppConfig struct {
	ApiHost 		string  	`json:"api_host"`    //api请求host
	IMSocketHost	string		`json:"im_socket_host"` //IM通讯host
	IMSocketPort    int			`json:"im_socket_port"` //IM通讯 port
	DomainName		string		`json:"domain_name"`  //域名
	Version			string		`json:"version"`	//版本
}


/*
读取配置文件
*/
func ReadConfig(path string) error {
	config = new(IMConfig)
	err := config.Parse(path)
	return err
}


/*
解析配置文件
*/
func (this *IMConfig) Parse(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &this)
	if err != nil {
		return err
	}
	return nil
}

/*
连接数据库
*/
func (this *DBConfig) Connect() (*sql.DB, error) {
	// 从配置文件中读取配置信息并初始化连接池(go中含有连接池处理机制)
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8", this.Username, this.Password, this.Host, this.Name)
	db, err := sql.Open("mysql", url)
	db.SetMaxIdleConns(this.MaxIdleConns) // 最大空闲连接
	db.SetMaxOpenConns(this.MaxOpenConns) // 最大连接数
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
