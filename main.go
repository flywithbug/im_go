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
	log.SetLevel(log.DEBUG)
	log.SetLayout("2006-01-02 15:04:05")
}

func main() {
	SetLog()
	defer log.Close()

	log.Info("*********************************************")
	log.Info("           系统:[%s]版本:[%s]", Name, Version)
	log.Info("*********************************************")
	configPath := flag.String("config", "config.json", "Configuration file to use")

	flag.Parse()

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal("读取配置文件错误:", err.Error())
	}
	model.Config = conf

	model.Database, err = conf.DBConfig.Connect()
	defer model.Database.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	go func() {
		err := server.StartHttpServer(*conf)
		log.Fatal("Http Server", err)
	}()

	perf.Init(conf.PprofBind)
	im.Listen(conf.IMPort)
}
