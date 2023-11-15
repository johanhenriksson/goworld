package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/johanhenriksson/goworld/game/net"
	"github.com/johanhenriksson/goworld/math/vec3"

	"github.com/quic-go/quic-go"
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
	sck, err := quic.ListenAddr(
		fmt.Sprintf("127.0.0.1:%d", net.GamePort),
		generateTLSConfig(),
		&quic.Config{
			KeepAlivePeriod: 3 * time.Second,
			EnableDatagrams: true,
		},
	)
	if err != nil {
		return err
	}

	log.Println("server: listening on port", net.GamePort)
	go func() {
		for {
			conn, err := sck.Accept(context.Background())
			if err != nil {
				panic(err)
			}
			go s.accept(conn)
		}
	}()
	return nil
}

func (s *Server) accept(conn quic.Connection) {
	log.Println("server: accepted client")
	client, err := NewClient(conn)
	if err != nil {
		client.Drop("failed to create client")
		return
	}

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

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"goworld"},
	}
}
