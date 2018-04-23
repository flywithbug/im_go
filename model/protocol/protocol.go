package protocol

import (
	"strconv"
	"regexp"
	"errors"
)

const  (
	protocol_header = "www.flywithbug.com:"
	boundary = "-----flywithbug-----"
	bufferLength = 1024
)

var lengthRe = regexp.MustCompile(`(www.flywithbug.com:)([0-9]+)(-----flywithbug-----)`)

func Packet(message []byte)[]byte  {
	totalByts := []byte(protocol_header)
	totalByts = append(totalByts,[]byte(strconv.Itoa(len(message)))...)
	totalByts = append(totalByts,[]byte(boundary)...)
	totalByts = append(totalByts,message...)
	return totalByts
}

func GetHeader(b []byte)([][]byte,error)  {
	match := lengthRe.FindSubmatch(b)
	if len(match) == 4	 {
		return match,nil
	}
	return nil,errors.New("数据校验失败")
}




