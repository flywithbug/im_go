package file

import (
	"errors"
	"os"
	"path"
	"io"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"io/ioutil"
)


// FileExists return a bool value that file is exist or not
func FileExists(path string) bool {
	if len(path) == 0 {
		return false
	}
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !stat.IsDir()
}


// CopyFile copy file from src to des
func CopyFile(src, des string, mkdir bool) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	if mkdir {
		dir := path.Dir(des)
		if !DirectoryExists(dir) {
			if err := Mkdir(dir, os.ModePerm); err != nil {
				return 0, err
			}
		}
	}

	desFile, err := os.Create(des)
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

// ReadLine is a function to read file which specify by filePath
func ReadLine(filePath string, callback ReadLineCallbackFunc) {
	if callback == nil {
		return
	}

	f, err := os.Open(filePath)
	defer f.Close()

	stop := false

	if err != nil {
		callback("", true, err, &stop)
		return
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				callback("", true, nil, &stop)
				break
			} else {
				callback("", false, err, &stop)
			}
		} else {
			callback(string(b), false, nil, &stop)
		}

		if stop {
			break
		}
	}
}

// FileMD5 return file md5
func FileMD5(file *os.File) (string, error) {
	md5Ctx := md5.New()
	_, err := io.Copy(md5Ctx, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5Ctx.Sum(nil)), nil
}


// WriteFile file directory decursively
func WriteFile(filePath string, data []byte, cover bool, mode os.FileMode) error {
	if data == nil {
		data = make([]byte, 0, 0)
	}
	filePath = strings.TrimSpace(filePath)
	if filePath == "" || !path.IsAbs(filePath) {
		return errors.New("请正确设置文件的绝对路径！")
	}
	if !cover && FileExists(filePath) {
		return errors.New("文件存在 [" + filePath + "]")
	}
	e := Mkdir(path.Dir(filePath), mode)
	if e != nil {
		return e
	}
	e = ioutil.WriteFile(filePath, data, mode)
	if e != nil {
		return e
	}
	return nil
}