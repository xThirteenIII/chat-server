package server

import (
	"chatserver/pkg/client"
	"fmt"
	"net"
	"sync"
)

func BroadcastMessages(broadcastChannel <-chan string, clients map[net.Conn]client.Client, mutex *sync.Mutex){

    // continously listen for strings in the channel
    // and send each of them to every client in the clients map
    for {
        msg := <- broadcastChannel

        // Prevent simoultaneous access to the clients map
        // data races avoided
        mutex.Lock()
        for _, client := range clients {
            fmt.Fprintln(client.Connection, msg)
        }
        mutex.Unlock()
    }

}
