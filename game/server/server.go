package server

import (
	"fmt"
	"log"
	osnet "net"

	"github.com/johanhenriksson/goworld/game/net"
)

type AuthToken struct{}

type Server struct {
	auths    map[uint64]AuthToken
	Instance *Instance
}

func NewServer() (*Server, error) {
	server := &Server{
		auths: map[uint64]AuthToken{
			1337: {},
		},
		Instance: NewInstance(),
	}
	if err := server.Listen(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) Listen() error {
	sck, err := osnet.Listen("tcp", fmt.Sprintf(":%d", net.GamePort))
	if err != nil {
		return err
	}
	log.Println("server: listening on port", net.GamePort)
	go func() {
		for {
			conn, err := sck.Accept()
			if err != nil {
				panic(err)
			}
			go s.accept(conn)
		}
	}()
	return nil
}

func (s *Server) accept(conn osnet.Conn) {
	log.Println("server: accepted client")
	client := NewClient(conn)

	pkt, err := client.decode()
	if err != nil {
		client.Drop("failed to decode auth packet")
		return
	}
	defer pkt.Message().Release()

	if pkt.Which() != net.Packet_Which_auth {
		client.Drop("expected auth packet")
		return
	}

	auth, err := pkt.Auth()
	if err != nil {
		client.Drop("invalid auth packet")
		return
	}
	log.Println("server: client auth with token", auth.Token())

	// authenticate client
	_, authed := s.auths[auth.Token()]
	if !authed {
		client.Drop("invalid authenticaton token")
		return
	}

	// invalidate authentication token
	delete(s.auths, auth.Token())

	// create player entity
	// todo: load player information
	player := &Player{
		id: 131273,
	}

	// enter world
	s.Instance.SubmitEvent(&EnterWorldEvent{
		Client: client,
		Player: player,
	})

	// start client read loop
	go client.readLoop()
}
