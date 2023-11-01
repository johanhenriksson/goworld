package server

import (
	osnet "net"

	"github.com/johanhenriksson/goworld/game/net"
)

type Server struct {
	auths map[uint64]Entity
}

func (s *Server) Listen() error {
	sck, err := osnet.Listen("tcp", ":420")
	if err != nil {
		return err
	}
	for {
		conn, err := sck.Accept()
		if err != nil {
			return err
		}
		go s.accept(conn)
	}
}

func (s *Server) accept(conn osnet.Conn) {
	client := NewClient(conn)

	pkt, err := client.decode()
	if err != nil {
		client.Drop("failed to decode auth packet")
		return
	}
	defer pkt.Message().Release()

	auth, ok := pkt.(*net.AuthPacket)
	if !ok {
		client.Drop("expected auth packet")
		return
	}

	// authenticate client
	entity, authed := s.auths[auth.Token()]
	if !authed {
		client.Drop("invalid authenticaton token")
		return
	}

	// invalidate authentication token
	delete(s.auths, auth.Token())

	// assign client to instance

	// trigger observe entity
	if err := client.Observe(entity); err != nil {
		client.Drop("faied to observe entity")
		return
	}

	// start client read loop
	go client.readLoop()
}
