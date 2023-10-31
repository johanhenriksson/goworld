package server

import (
	"log"

	"capnproto.org/go/capnp/v3"
	"github.com/johanhenriksson/goworld/game/net"
)

type Client struct {
	Entity
	encoder capnp.Encoder
	decoder capnp.Decoder
}

func (c *Client) readLoop() {

	// read network packets and submit events to the current instance
	for {
		msg, err := c.decoder.Decode()
		if err != nil {
			log.Fatalln(err)
		}
		// todo: how to figure out what kind of packet it is?
		// todo: submit packet event to instance
		msg.Release()
	}
}

func (c *Client) Send(packet net.Packet) error {
	return c.encoder.Encode(packet.Message())
}
