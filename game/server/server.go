package server

import (
	"log"
	"net"

	packet "github.com/johanhenriksson/goworld/game/net"

	"capnproto.org/go/capnp/v3"
)

type Server struct {
}

func (s *Server) Listen() error {
	sck, err := net.Listen("tcp", ":420")
	if err != nil {
		return err
	}
	for {
		conn, err := sck.Accept()
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	decoder := capnp.NewDecoder(conn)

	// expect the first message to be an authentication packet
	msg, err := decoder.Decode()
	if err != nil {
		log.Println("accept failed:", err)
	}
	auth, err := packet.ReadRootAuthPacket(msg)
	if err != nil {
		log.Println("accept failed:", err)
	}

	auth.Token()
}
