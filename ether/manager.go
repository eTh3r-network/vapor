package ether

import "net"

type Manager struct {
    listenPort  int
    stop        bool
}

func Initialise(port int) (*Manager) {
    newManager := new(Manager)
    newManager.listenPort = port

    return Manager
}

func (*Manager) HandleConnection(conn net.Conn) (error) {
    // Do stuff
}

func (*Manager) Listen() (error) {
    ln, err := net.Listen("tcp", ":"+listenPort)

    if err != nil {
        logger.Warn("An error has been encountered while trying to bind to the port: ", err)
        return err
    }

    for !stop {
        conn, err := ln.Accept()

        if err != nil {
            logger.Warn("An error has been encountered while a user is connecting")
        }

        logger.Info("Handling connection from", conn.RemoteAddr().String())
        go HandleConnection(conn)
    }
}
