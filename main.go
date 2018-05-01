package main

import (
	"flag"
	"im_go/config"
	"im_go/libs/perf"
	"im_go/model"
	"im_go/server"
	"log"
	//"im_go/ims"
	"im_go/im"
)

const (
	Name    string = "IM"
	Version string = "1.0"
)



func main() {
	log.Println("*********************************************")
	log.Printf("           系统:[%s]版本:[%s]", Name, Version)
	log.Println("*********************************************")
	configPath := flag.String("config", "config.json", "Configuration file to use")
	flag.Parse()

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("读取配置文件错误: %s", err)
	}
	model.Config = conf

	model.Database, err = conf.DBConfig.Connect()
	defer model.Database.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}
	go func() {
		err := server.StartHttpServer(*conf)
		log.Fatalln("Http Server", err)
	}()

	perf.Init(conf.PprofBind)
	//初始化model包下全局变量值
	//ims.Listen(conf.IMPort)
	im.Listen(conf.IMPort)
}
