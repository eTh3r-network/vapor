package ether

import "encoding/binary"

func (c *Connection) serve0001() {
	c.authState = 1
	
	if err := c.ack(); err != nil {
		c.abandon()
		return
	}

	var buff []byte
	pass := false

	for !pass {
		l, err := c.bind.Read(buff)

		// Check minimal values
		if err != nil || l < 4 {
			c.log.Warn("There has been an error while receiving the key", err)
			c.handleErr(1, 0xaa)

			continue
		}

		// Check pkt constant
		if buff[0] != 0x0e || buff[1] != 0x1f {
			c.log.Warn("Key packet malformed")
			c.handleErr(1, 0xab)

			continue
		}

		// Extract key length
		keyLength := binary.BigEndian.Uint16(buff[2:4])

		// Verify message length
		if uint16(len(buff)) != uint16(4) + keyLength {
			c.log.Warn("Key payload malformation")
			c.handleErr(1, 0xac)

			continue
		}

		// Store key and key length
		c.keyLength = keyLength 
		c.key = buff[4:4+keyLength]

		pass = true
	}

	c.authState = 2

	if err := c.ack(); err != nil {
		c.abandon() // there is no reason to reach this
		return
	}

}

