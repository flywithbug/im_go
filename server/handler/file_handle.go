package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"path/filepath"
	"os"
	"io"
	"github.com/pborman/uuid"
	"github.com/flywithbug/file"
	"github.com/flywithbug/log4go"
	"time"
	"im_go/model"
)

const  localFilePath  =  "./image/"

func UploadImageHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()

	//gin将het/http包的FormFile函数封装到c.Request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}

	if header.Size > 1024*170 {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusRequestEntityTooLarge,fmt.Sprintf(" file to big no more than 150kb "))
		return
	}
	today := time.Now().Format("2006-01-02")
	localpath := localFilePath+today+"/"
	//获取文件名
	ext := filepath.Ext(header.Filename)
	name:= uuid.New()
	filename := name + ext
	//写入文件
	bexit,err := PathExists(localpath)
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	if !bexit {
		err = os.Mkdir(localpath, os.ModeDir)
		if err != nil {
			log4go.Info(err.Error())
			aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
			return
		}
	}

	out, err := os.Create(localpath + filename)
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	avatarpath  := fmt.Sprintf("filepath=%s&dir=%s",filename,today)
	user , _ := User(c)
	if user != nil{
		model.UpdateUserAvatar(avatarpath,user.UserId)
	}
	aRes.SetResponseDataInfo("filepath",avatarpath)
}

func DownloadImageHandler(c *gin.Context)  {
	filepath := c.Query("filepath")
	dir := c.Query("dir")
	log4go.Info(dir,filepath)
	filename :=  localFilePath + dir +"/"+filepath
	if !file.FileExists(filename) {
		c.Writer.Write([]byte("Error: Image Not found."))
		return
	}
	http.ServeFile(c.Writer,c.Request,filename)
}


// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

