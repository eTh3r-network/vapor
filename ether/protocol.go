package ether

import "fmt"

type Room struct {

}

type Connection struct {
    authState   int
    key         int[]
    keyId       int[]
    keyLength   uint16
    rooms       *Room[]
}

func Test() {
    fmt.Println("Ola from protocol")
}
