//
// This file is part of the eTh3r project, written, hosted and distributed under MIT License
//  - eTh3r network, 2023-2024
//

package ether

import "fmt"
import "net"
import "log/slog"
import "encoding/binary"

type Room struct {
	roomId	uint8
}

type Connection struct {
	authState   	int
	key         	[]byte
	keyId       	[]byte
	keyLength   	uint16
	rooms       	[]*Room

	bind		net.Conn
	log 		*slog.Logger
}

func InitialiseConnection(conn net.Conn, log *slog.Logger) (*Connection) {
	newConn := new(Connection)
	
	newConn.authState = 0
	newConn.bind = conn
	newConn.log = log

	return newConn
}

func (c *Connection) handleErr(level int, _type byte) {
	_, err := c.bind.Write([]byte{_type})

	if err != nil {
		c.log.Warn("There has been an error transmitting the error code. Bad luck .-.", "level", level)
	}
}

func (c *Connection) ack() (error) {
	_, err := c.bind.Write([]byte{0x0a})

	if err != nil {
		c.log.Warn("There has been an error while sending ack packet to client:", err)
	}

	return err
}

func (c *Connection) abandon() {
	_ = c.bind.Close()
}

func (c *Connection) Serve() {
	var buff []byte

	l, err := c.bind.Read(buff)

	if err != nil {
		c.log.Warn("An error has been caught reading a message")
		c.handleErr(0, 0xff)
	}

	if l != 6 {
		c.log.Warn("The client sent a wrong amount of data")
		c.handleErr(0, 0xa1)
	}

	if buff[0] != 0x05 || buff[1] != 0x31 || buff[2] != 0x80 || buff[3] != 0x08 {
		c.log.Warn("Wrong payload, first message")
		c.handleErr(0, 0xa2)
	}

	version := binary.BigEndian.Uint16(buff[4:6])
	
	switch version {
		case 0x0001:
			c.serve0001()
		default:
			c.log.Warn("Unsupported version", "ver", version)
			c.handleErr(0, 0xa4)
	}
}

func Test() {
	fmt.Println("Ola from protocol")
}
