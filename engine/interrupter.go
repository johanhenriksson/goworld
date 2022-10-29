package engine

import (
	"log"
	"os"
	"os/signal"
)

type Interrupter interface {
	Running() bool
}

type interrupter struct {
	running bool
}

func (r *interrupter) Running() bool {
	return r.running
}

func NewInterrupter() Interrupter {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	r := &interrupter{running: true}

	go func() {
		for range sigint {
			if !r.running {
				log.Println("Kill")
				os.Exit(1)
			} else {
				log.Println("Interrupt")
				r.running = false
			}
		}
	}()

	return r
}
