package server

import (
	"fmt"
	"log"
	osnet "net"

	"github.com/johanhenriksson/goworld/game/net"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type AuthToken struct {
	Token uint64
}

type Server struct {
	Instance *Instance
}

func NewServer() (*Server, error) {
	server := &Server{
		Instance: NewInstance(),
	}
	if err := server.Listen(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) Authenticate(token uint64) (*Player, error) {
	return &Player{
		id:       Identity(token),
		position: vec3.New(0, 10, 0),
	}, nil
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
	player, err := s.Authenticate(auth.Token())
	if err != nil {
		client.Drop("invalid authenticaton token")
		return
	}

	// enter world
	s.Instance.SubmitEvent(&EnterWorldEvent{
		Client: client,
		Player: player,
	})

	// start client read loop
	go client.readLoop()
}
