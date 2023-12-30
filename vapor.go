//
// This file is part of the eTh3r project, written, hosted and distributed under MIT License
//  - eTh3r network, 2023-2024
//

package main

import "fmt"

import "github.com/eTh3r-network/vapor/ether"
import "github.com/eTh3r-network/vapor/logger"


func main() {
    fmt.Println("ola :D")
    ether.Test()

    log := logger.GetLogger()

    manager := ether.Initialise(2142, log)
    _ = manager.Listen()
}
