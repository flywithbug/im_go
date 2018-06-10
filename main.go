package main

import (
	"flag"
	"im_go/config"
	"im_go/libs/perf"
	"im_go/model"
	"im_go/server"

	log "github.com/flywithbug/log4go"
	"im_go/im"
)

const (
	Name    string = "IM"
	Version string = "1.0"
)


func SetLog() {
	w := log.NewFileWriter()
	w.SetPathPattern("./log/log-%Y%M%D.log")
	c := log.NewConsoleWriter()
	c.SetColor(true)
	log.Register(w)
	log.Register(c)
	log.SetLevel(config.Conf().LogLevel%4)
	log.SetLayout("2006-01-02 15:04:05")
}

func main() {
	configPath := flag.String("config", "config.json", "Configuration file to use")
	flag.Parse()
	//加载配置文件
	err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal("读取配置文件错误:", err.Error())
	}
	SetLog()
	defer log.Close()

	//连接数据库
	model.Database, err= config.Conf().DBConfig.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer model.Database.Close()

	go func() {
		//启动用户管理服务
		server.StartHttpServer(config.Conf().ServerPort,config.Conf().RouterPrefix)
	}()

	//启动系统监控
	perf.Init(config.Conf().ProfBind)

	//启用im服务
	im.StartIMServer(config.Conf().IMPort,config.Conf().HttpPort)
}
