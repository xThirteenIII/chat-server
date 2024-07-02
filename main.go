package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Conn is a generic stream-oriented network connection.
//
// Multiple goroutines may invoke methods on a Conn simultaneously.
type Client struct {
    conn net.Conn
    name string
}

var (

    // Clients associates each connection to a Client object
    clients = make(map[net.Conn]Client)

    // Broadcast is a channel for broadcasting messages to all channels
    broadcast = make(chan string)

    // A mutex to ensure thread-safe access to the clients map.
    mutex = &sync.Mutex{}
)

// handleConnection manages each client’s connection, reads their messages, and broadcasts them.
func handleConnection(connection net.Conn){

    // Close connection when finished
    defer connection.Close()

    // NewReader returns a new Reader whose buffer has the default size.
    reader := bufio.NewReader(connection)
    fmt.Fprint(connection, "Enter your name: " )
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    client := Client{conn: connection, name: name}
    mutex.Lock()
    clients[connection] = client
    mutex.Unlock()

    broadcast <- fmt.Sprintf("%s has joined the chat", name)

    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        broadcast <- fmt.Sprintf("%s: %s", name, strings.TrimSpace(message))
    }

    mutex.Lock()

    // The delete built-in function deletes the element with the specified key
    // (m[key]) from the map. If m is nil or there is no such element, delete
    // is a no-op.
    delete(clients, connection)
    mutex.Unlock()
    broadcast <- fmt.Sprintf("%s has left the chat", name)
}

// handleMessages listens for messages on the broadcast channel and sends them to all connected clients.
func handleMessages(){
    for {
        msg := <- broadcast
        mutex.Lock()
        for _, client := range clients {
            fmt.Fprintln(client.conn, msg)
        }
        mutex.Unlock()
    }
}


func main(){
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server: ", err)
        return
    }
    defer listener.Close()

    go handleMessages()

    fmt.Println("Chat server started on :8080")
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection", err)
            continue
        }
        go handleConnection(conn)
    }
}
