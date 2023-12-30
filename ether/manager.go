package ether


import "net"
import "log/slog"
import "strconv"

type Manager struct {
	listenPort  	int
	stop		bool
	logger		*slog.Logger
	clients		[]*Connection
	rooms		[]*Room
}

func Initialise(port int, log *slog.Logger) (*Manager) {
	newManager := new(Manager)

	newManager.listenPort = port
	newManager.stop = false
	newManager.logger = log

	return newManager
}

func (m *Manager) Listen() (error) {
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(m.listenPort))
	
	if err != nil {
		m.logger.Warn("An error got caught while trying to bind:", err)
		return err
	}

	for !m.stop {
		conn, err := ln.Accept()

		if err != nil {
			m.logger.Warn("There has been an issue when a client tried to connect:", err)
		}

		m.logger.Info("Serving", conn.RemoteAddr().String())

		client := InitialiseConnection(conn, m.logger)
		m.clients = append(m.clients, client)

		go client.Serve()
	}

	return nil
}
