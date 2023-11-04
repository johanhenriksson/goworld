package main

import (
	"log"
	"time"

	"github.com/johanhenriksson/goworld/game/server"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("server instance:", srv.Instance)

	for {
		time.Sleep(1 * time.Second)
	}
}
