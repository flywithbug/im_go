package model

import (
	"database/sql"
	"im_go/config"
)

/*
 包内上下文变量
 */
var (
	Database *sql.DB        = nil //数据库操作对象
	Config   *config.IMConfig
)
