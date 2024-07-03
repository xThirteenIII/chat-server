package server

import (
	"bufio"
	"chatserver/pkg/client"
	"fmt"
	"net"
	"strings"
	"sync"
)

func HandleConnection(connection net.Conn, clients map[net.Conn]client.Client, broadcast chan <- string, mutex *sync.Mutex, waitgroup *sync.WaitGroup){

    defer waitgroup.Done()

    // Close connection when function ends
    defer connection.Close() 

    // create new reader
    reader := bufio.NewReader(connection)

    // ask for user name
    fmt.Fprintf(connection, "Enter your name: ")

    // save user name
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    // create unique client based on current connection and name
    client := client.Client{Connection: connection, Name: name}

    // add client to the map
    mutex.Lock()
    clients[connection] = client
    mutex.Unlock()

    // send the name to the broadcast channel
    broadcast <- fmt.Sprintf("%s has joined the chat", name)

    // continously wait for user messages and sent them to the channel
    for {
        message, err := reader.ReadString('\n')
        if err!=nil{
            break
        }
        broadcast <- fmt.Sprintf("%s: %s",name, strings.TrimSpace(message))
               
    }

    // delete client corresponding to current connection from the map
    mutex.Lock()
    delete(clients, connection)
    mutex.Unlock()

    // send leave string to the channel
    broadcast <- fmt.Sprintf("%s has left the chat", name)
}
