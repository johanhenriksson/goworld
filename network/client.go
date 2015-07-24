package network

import (
    "net"
)

const (
    HeaderLength     = 8
    ClientBufferSize = 1024
    ClientBacklog    = 10
)

type Header struct {
    Protocol        uint32
    SeqNum          uint16
    Ack             uint16
    AckField        uint32
}

type Packet struct {
    Header
}

type Client struct {
    *Socket
}

func NewClient(conn *net.UDPConn, addr *net.UDPAddr) *Client {
    return &Client {
        Socket: NewSocket(conn, addr, ClientBufferSize, false),
    }
}

func NewServerClient(conn *net.UDPConn, addr *net.UDPAddr) *Client {
    return &Client {
        Socket: NewSocket(conn, addr, ClientBufferSize, true),
    }
}

func ConnectTo(hostname string) (*Client, error) {
    srv_addr, err := net.ResolveUDPAddr("udp", hostname)
    if err != nil { return nil, err }

    conn, err := net.DialUDP("udp", nil, srv_addr)
    if err != nil { return nil, err }

    return NewClient(conn, srv_addr), nil
}

func (c *Client) Worker() {
    for {
        _, _, err := c.Read(c.buffer)
        if err != nil { panic(err) }

        /* TODO: Check protocol id */
        c.Recv(c.buffer)
    }
}
