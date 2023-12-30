package ether

import "fmt"
import "net"

type Room struct {
	roomId	uint8
}

type Connection struct {
	authState   	int
	key         	[]int
	keyId       	[]int
	keyLength   	uint16
	rooms       	[]*Room

	bind		net.Conn
}

func InitialiseConnection(conn net.Conn) (*Connection) {
	newConn := new(Connection)
	
	newConn.authState = 0
	newConn.bind = conn

	return newConn
}

func (*Connection) Serve() {
	
}


func Test() {
	fmt.Println("Ola from protocol")
}
