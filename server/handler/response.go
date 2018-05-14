package handler

import "net/http"

//type Response map[string]interface{}

type Response struct {
	Data 	map[string]interface{}	`json:"data"`
	Code 	int						`json:"code"`
}


func NewResponse() *Response {
	res := new(Response)
	res.Code = http.StatusOK
	res.Data = make(map[string]interface{})
	return res
}


func ErrorResponse(errno int, msg string) *Response {
	r := NewResponse()
	r.Code = errno
	r.SetErrorInfo(errno, msg)
	return r
}

func (s *Response) SetErrorInfo(errno int, msg string) {
	s.Code = errno
	s.Data["msg"] = msg
}
func (s *Response) SetSuccessInfo(code int, msg string) {
	s.Code = code
	s.Data["msg"] = msg
}

func (s *Response) AddResponseInfo(key string, val interface{}) {
	s.Data[key] = val
}


