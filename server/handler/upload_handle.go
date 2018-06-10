package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"path/filepath"
	"os"
	"io"
	"github.com/pborman/uuid"
)


func UploadFile(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(aRes.Code,aRes)
	}()

	//gin将het/http包的FormFile函数封装到c.Request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	//获取文件名
	ext := filepath.Ext(header.Filename)
	name:= uuid.New()
	filename := name + ext
	//写入文件
	out, err := os.Create("./file/" + filename)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}
	aRes.SetResponseDataInfo("filepath",filename)

}