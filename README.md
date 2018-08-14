
//项目说明：存储部分使用的mysql, 暂未增加redis支持
//TODO: Redis 支持。


##流程图


## 安装说明

//编译为centos 执行文件
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build main.go


系统需要安装Go和MySQL。


打开配置文件 config.json，修改相关配置。


创建数据库im_go，再导入install.sql

	$ mysql -u username -p -D im_go < install.sql




