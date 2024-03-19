//
// This file is part of the eTh3r project, written, hosted and distributed under MIT License
//  - eTh3r network, 2023-2024
//

package ether

import "encoding/binary"

func (c *Connection) serve0001(manager *Manager) {
	c.authState = 1

	if err := c.ack(); err != nil {
		c.abandon()
		return
	}

	var buff []byte
	pass := false

	for !pass {
		l, err := c.bind.Read(buff)

		if err != nil {
			c.log.Warn("There has been an error in pkt reading", err)
			c.handleErr(1, 0xff)

			continue
		}

		// Check minimal length
		if l < 4 {
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
		if uint16(len(buff)) != uint16(4)+keyLength {
			c.log.Warn("Key payload malformation")
			c.handleErr(1, 0xac)

			continue
		}

		// Store key and key length
		c.keyLength = keyLength
		c.key = buff[4 : 4+keyLength]

		pass = true
	}

	c.authState = 2

	if err := c.ComputeKeyId(); err != 0 {
		c.log.Warn("Could not derive KeyId from PubKey")
		c.handleErr(2, 0xad)
	}

	manager.RegisterConnection(c)

	c.authState = 3

	if err := c.ack(); err != nil {
		c.abandon() // there is no reason to reach this
		return
	}

	for {
		l, err := c.bind.Read(buff)

		if err != nil {
			c.log.Warn("There has been an error reading pkt", err)
			c.handleErr(3, 0xff)

			continue
		}

		// should at least have a cons
		if l < 1 {
			c.log.Warn("Packet malformed")
			c.handleErr(3, 0xba)
		}

		switch buff[0] {
		case 0xba:
			c.log.Debug("Fetching user", buff[1:])
			conn := manager.FetchUserById(buff[1:])

			if conn == nil {
				c.log.Warn("Could not find user")

				respBuff := []byte{0xca}
				respBuff = append(respBuff[:], buff[:]...)

				_, err := c.bind.Write(respBuff)

				if err != nil {
					c.log.Warn("Could not send message")
				}

				continue
			}

			keyBuff := make([]byte, 2)
			binary.LittleEndian.PutUint16(keyBuff, conn.keyLength) // append the key length as uint16
			keyBuff = append(keyBuff, conn.key...)                 // append the key

			respBuff := []byte{0xa0, 0xba}
			respBuff = append(respBuff[:], keyBuff[:]...) // prepend the pck id

			_, err := c.bind.Write(respBuff) // write

			if err != nil {
				c.log.Warn("Could not send the user key, internal server error")
			}

			continue
		case 0xee:
			// Knock
			c2 := buff[1:]

			var kLength uint16 = 0
			binary.LittleEndian.PutUint16(c2[:2], kLength)

			if uint16(len(c2)-4) != kLength {
				c.log.Warn("Wrong packet length")
				c.handleErr(2, 0xa1)

				continue
			}

			c.log.Debug("Fetching user", c2[2:])
			c2Conn := manager.FetchUserById(c2[2:])

			if c2Conn == nil {
				c.log.Warn("Could not find c2")
				c.handleErr(2, 0xad) // TODO: add 0xad as user not found

				continue
			}

			c2Conn.SendKnock0001(c)

			if _, err := c.bind.Write([]byte{0xa0, 0xee}); err != nil {
				c.log.Warn("Could not send the ack pkg")

				continue
			}

			continue
		case 0xab:
			// Knock ans

			respVal := (buff[1] == 0x01)
			c.log.Debug("%b", respVal)

			c2 := buff[2:]

			var kLength uint16 = 0
			binary.LittleEndian.PutUint16(c2[:2], kLength)

			if uint16(len(c2)-5) != kLength {
				c.log.Warn("Wrong packet length")
				c.handleErr(2, 0xa1)

				continue
			}

			c.log.Debug("Fetching user", c2[2:])
			c2Conn := manager.FetchUserById(c2[2:])

			if c2Conn == nil {
				c.log.Warn("Could not find c2")
				c.handleErr(2, 0xad) // TODO: add 0xad as user not found

				continue
			}

			if _, err := c2Conn.bind.Write(buff); err != nil {
				c2Conn.log.Warn("Could not send pkg")
			}

			r := manager.SpawnRoom()

			c.NotifyRoom(r, c2Conn)
			c2Conn.NotifyRoom(r, c)

			continue
		case 0xda:
			// Message
			var ridLength uint8 = buff[1]
			rid := buff[2 : 2+ridLength]

			room := manager.FetchRoom(rid)
			room.SendMessageToRecipients0001(buff, c)

			if _, err := c.bind.Write([]byte{0xa0, 0xda}); err != nil {
				c.log.Warn("Could not send ack to message")
			}

			continue
		case 0xaf:
			// Room termination

			var ridLength uint8 = buff[1]
			rid := buff[2 : 2+ridLength]

			room := manager.FetchRoom(rid)
			room.SendMessageToAllRecipients0001(buff)

			manager.TerminateRoom(room)

			continue
		case 0xbf:
			// Disconnect
			c.log.Info("Client disconnecting")
			_ = c.ack()

			return

		default:
			c.log.Warn("Unknown packet:", "pkid", buff[0])
			c.handleErr(2, 0xfd)
			continue
		}
	}
}

func (c *Connection) SendKnock0001(c2 *Connection) {
	return
}

func (c *Connection) NotifyRoom(r *Room, c2 *Connection) {
	return
}

func (r *Room) SendMessageToAllRecipients0001(buff []byte) []error {
	var errs []error

	return errs
}

func (r *Room) SendMessageToRecipients0001(buff []byte, sender *Connection) []error {
	var errs []error

	return errs
}
