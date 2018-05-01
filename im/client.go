package im

import (
	"net"
	"log"
	"bufio"
)

type Client struct {
	Connection

}

func NewClient(conn *net.TCPConn)*Client  {
	client := new(Client)
	client.conn = conn

	return client
}


func (client *Client)Read()  {

}


func (client *Client)Write()  {

}


func (client *Client) Listen() {
	go client.Read()
	go client.Write()

}
