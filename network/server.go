package network

import (
    "fmt"
    "net"
)

const (
    ServerBufferSize = 1024
)

type Server struct {
    *Socket
    clients     map[string]*Client
}

func NewServer(conn *net.UDPConn, addr *net.UDPAddr) *Server {
    return &Server {
        Socket: NewSocket(conn, addr, ServerBufferSize, true),
        clients: make(map[string]*Client),
    }
}

func Listen(endpoint string) (*Server, error) {
    addr, err := net.ResolveUDPAddr("udp", endpoint)
    if err != nil {
        return nil, err
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        return nil, err
    }

    return NewServer(conn, addr), nil
}

func (s *Server) Worker() {
    for {
        _, in_addr, err := s.Read(s.buffer)
        if err != nil {
            panic(err)
        }
        str_addr := in_addr.String()

        /* TODO: Check protocol id */

        var client *Client
        if c, ok := s.clients[str_addr]; ok {
            client = c
        } else {
            client = NewServerClient(s.conn, in_addr)
            s.clients[str_addr] = client
            fmt.Println("Client from", in_addr)
        }

        client.Recv(s.buffer)
        client.Send([]byte("REPLY!"))
    }
}

func (s *Server) Update(dt float64) {
    for _, client := range s.clients {
        client.Update(dt)
    }
}

func (s *Server) Stop() {
    s.conn.Close()
}
