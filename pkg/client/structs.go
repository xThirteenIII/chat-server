package client

import "net"

type Client struct{
    Connection net.Conn
    Name string
}
