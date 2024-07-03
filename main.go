package main

import (
	"chatserver/pkg/client"
	"chatserver/pkg/server"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
    clients = make(map[net.Conn]client.Client)
    broadcast = make(chan(string), 100)
    mutex = &sync.Mutex{}
    waitgroup  sync.WaitGroup
    address = ":8080"
    network = "tcp"
)

func main(){

    // Listen announces on the local network address.
    listener, err := net.Listen(network, address)
    if err != nil {
        os.Exit(1)
    }
    defer listener.Close()

    go server.BroadcastMessages(broadcast, clients, mutex)

    fmt.Println("Chat server started on :8080")

    // continously listen for connections
    for {
        
	    // Accept waits for and returns the next connection to the listener.
        connection, err := listener.Accept()
        if err != nil {
            fmt.Println("Failed accepting connection")
            continue
        }

        waitgroup.Add(1)

        go server.HandleConnection(connection, clients, broadcast, mutex, &waitgroup)
    }

    waitgroup.Wait()
}
