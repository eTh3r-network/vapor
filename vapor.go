package main

import "fmt"

import "github.com/eTh3r-network/vapor/ether"
import "github.com/eTh3r-network/vapor/logger"


func main() {
    fmt.Println("ola :D")
    ether.Test()

    log := logger.GetLogger()

    manager := ether.Initialise(2142, log)
    manager.Listen()
}
