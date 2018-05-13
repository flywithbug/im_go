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


func SetLog(conf *config.IMConfig) {
	w := log.NewFileWriter()
	w.SetPathPattern("./log/log-%Y%M%D.log")
	c := log.NewConsoleWriter()
	c.SetColor(true)
	log.Register(w)
	log.Register(c)
	log.SetLevel(conf.LogLevel%4)
	log.SetLayout("2006-01-02 15:04:05")
}

func main() {
	configPath := flag.String("config", "config.json", "Configuration file to use")
	flag.Parse()
	//加载配置文件
	conf, err := config.ReadConfig(*configPath)

	SetLog(conf)
	defer log.Close()

	log.Info("*********************************************")
	log.Info("           系统:[%s]版本:[%s]", Name, Version)
	log.Info("*********************************************")

	if err != nil {
		log.Fatal("读取配置文件错误:", err.Error())
	}


	//连接数据库
	model.Database, err= conf.DBConfig.Connect()
	defer model.Database.Close()
	if err != nil {
		log.Fatal(err.Error())
	}


	go func() {
		//启动用户管理服务
		server.StartHttpServer(conf.ServerPort,conf.RouterPrefix)
	}()

	//启动系统监控
	perf.Init(conf.PprofBind)

	//启用im服务
	im.StartIMServer(conf.IMPort,conf.HttpPort)
}
