package model

import "net"

/*
	抽象出每一个接入的client连接
*/
type Client struct {
	Key 		string       //conn 唯一标识
	Conn		net.Conn 	 //连接



}
