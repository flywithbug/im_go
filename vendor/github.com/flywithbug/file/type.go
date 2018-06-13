package file

import (
	"os"
	"time"
)

// ReadLineCallbackFunc is callback function type definition for ReadLine function
// line: One line of file
// finished: If is end of file, finished will be true
// err: Error for read line
// stop: If stop be setted true in callback function, the ReadLine function will stop
type ReadLineCallbackFunc func(line string, finished bool, err error, stop *bool)


// ListDirectCallbackFunc is callback function type definition for ListDirectory function
// file: file info object
// err: Error for list directory
type ListDirectCallbackFunc func(file FileInfo, err error)


// FileInfo redefine os.FileInfo
type FileInfo interface {
	os.FileInfo
	FilePath() string
}

// FileInfoBase implement FileInfo
type FileInfoBase struct {
	FileInfo os.FileInfo
	Path     string
}

// Name return file name
func (s *FileInfoBase) Name() string {
	return s.FileInfo.Name()
}

// Size return file size
func (s *FileInfoBase) Size() int64 {
	return s.FileInfo.Size()
}

// Mode return file mode
func (s *FileInfoBase) Mode() os.FileMode {
	return s.FileInfo.Mode()
}

// ModTime reutrn file last modify time
func (s *FileInfoBase) ModTime() time.Time {
	return s.FileInfo.ModTime()
}

// IsDir return the file is directory or not
func (s *FileInfoBase) IsDir() bool {
	return s.FileInfo.IsDir()
}

// Sys underlying data source (can return nil)
func (s *FileInfoBase) Sys() interface{} {
	return s.FileInfo.Sys()
}

// FilePath return file path
func (s *FileInfoBase) FilePath() string {
	return s.Path
}