package server

import (
	"chatserver/pkg/client"
	"fmt"
	"net"
	"sync"
)

func HandleMessages(broadcastChannel <-chan string, clients map[net.Conn]client.Client, mutex *sync.Mutex){

    // continously listen for strings in the channel
    for {
        msg := <- broadcastChannel
        mutex.Lock()
        for _, client := range clients {
            fmt.Println(client.Connection, msg)
        }
        mutex.Unlock()
    }

}
