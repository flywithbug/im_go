package file

import (
	"errors"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"path"
)

func DirectoryExists(path string)bool  {
	if len(path) == 0 {
		return false
	}
	fileinfo, err := os.Stat(path)
	if err != nil{
		return false
	}
	return fileinfo.IsDir()
}

func ListDirectory(path string, recursive bool,callback  ListDirectCallbackFunc)  {
	if callback == nil{
		return
	}
	if !DirectoryExists(path) {
		msg := "Directory path '" + path + "' is not exist!"
		callback(nil, errors.New(msg))
		return
	}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		callback(nil, err)
	}

	for _, fileinfo := range dir {
		fi := new(FileInfoBase)
		fi.FileInfo = fileinfo
		fi.Path = filepath.Join(path, fileinfo.Name())
		var newFileInfo FileInfo = fi
		callback(newFileInfo, nil)
		if recursive && newFileInfo.IsDir() {
			ListDirectory(newFileInfo.FilePath(), recursive, callback)
		}
	}
}


// CopyDirectory copy src directory to des
// if directory is not exists, create it using os.ModePerm
func CopyDirectory(src, des string) []error {
	if !DirectoryExists(src) {
		return []error{errors.New("Source directory is not exists")}
	}
	if !DirectoryExists(des) {
		if err := Mkdir(des, os.ModePerm); err != nil {
			return []error{err}
		}
	}

	errs := make([]error, 0, 5)
	ListDirectory(src, true, func(file FileInfo, err error) {
		if err != nil {
			errs = append(errs, err)
			return
		}
		ouputRel, err := filepath.Rel(src, file.FilePath())
		if err != nil {
			errs = append(errs, err)
			return
		}
		output := filepath.Join(des, ouputRel)
		if file.IsDir() {
			if err := Mkdir(output, os.ModePerm); err != nil {
				errs = append(errs, err)
			}
			return
		}
		if _, err := CopyFile(file.FilePath(), output, true); err != nil {
			errs = append(errs, err)
		}
	})
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Mkdir create directory decursively
func Mkdir(dir string, mode os.FileMode) error {
	dir = strings.TrimSpace(dir)
	if dir == "" || !path.IsAbs(dir) {
		return errors.New("请正确设置目录的绝对路径！")
	}
	if !DirectoryExists(dir) {
		e := Mkdir(path.Dir(dir), mode)
		if e != nil {
			return e
		}
		return os.Mkdir(dir, mode)
	}
	return nil
}
