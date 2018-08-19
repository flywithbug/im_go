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
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	"golang.org/x/image/bmp"
	"errors"
	"github.com/nfnt/resize"
	)

const  localFilePath  =  "./image/"

func UploadUserAvatarHandler(c *gin.Context)  {
	aRes := NewResponse()
	defer func() {
		c.JSON(http.StatusOK,aRes)
	}()

	//gin将het/http包的FormFile函数封装到c.Request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest,fmt.Sprintf("get file err : %s", err.Error()))
		return
	}

	//if header.Size > 1024*170 {
	//	aRes.SetErrorInfo(http.StatusRequestEntityTooLarge,fmt.Sprintf(" file to big no more than 150kb "))
	//	return
	//}
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
		err = os.Mkdir(localpath, os.ModePerm)
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
	size := c.Query("size")
	if len(size) == 0 {
		size = "120"
	}
	filename :=  localFilePath + dir +"/" +size +"-"+filepath
	fileOrigin :=  localFilePath + dir +"/" +filepath
	if !file.FileExists(filename)  {
		if !file.FileExists(fileOrigin) {
			c.Writer.Write([]byte("Error: Image Not found."))
			return
		}
		fIn, _ := os.Open(fileOrigin)
		log4go.Info(fileOrigin)
		defer fIn.Close()
		fOut, _ := os.Create(filename)
		log4go.Info(filename)
		defer fOut.Close()
		err := scale(fIn, fOut, 120, 120, 100)
		if err != nil {
			log4go.Info(err.Error())
			http.ServeFile(c.Writer,c.Request,fileOrigin)
			return
		}
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



func scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	//return jpeg.Encode(out, canvas, &jpeg.Options{quality})

	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, canvas)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}
