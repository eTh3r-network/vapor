package main

import "fmt"

import "github.com/eTh3r-network/vapor/ether"
import "github.com/eTh3r-network/vapor/logger"


func main() {
    fmt.Println("ola :D")
    ether.Test()

    manager := ether.Initialise()
    manager.Listen()
}
