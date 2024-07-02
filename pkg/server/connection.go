package server

import (
	"bufio"
	"chatserver/pkg/client"
	"net"
	"strings"
)

func HandleConnection(connection net.Conn, clients map[net.Conn]client.Client, broadcast chan <- string){
    defer connection.Close() 
    reader := bufio.NewReader(connection)
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
}
