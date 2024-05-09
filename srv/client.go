package srv

import "log"

type ClientToken struct {
	Character string
}

type Client interface {
	// Reads authentication token from client.
	// Blocks until the token is successfully read.
	// Returns an error if the next packet is not an authentication packet.
	ReadToken() (ClientToken, error)

	// Converts the event to a packet and sends it to the client.
	// Non-blocking.
	Send(Event)

	Drop(reason string) error
}

type DummyClient struct {
	Token ClientToken
}

var _ Client = (*DummyClient)(nil)

func (c *DummyClient) ReadToken() (ClientToken, error) {
	return c.Token, nil
}

func (c *DummyClient) Send(ev Event) {
	// do nothing
	log.Println(c.Token.Character, "<-", ev)
}

func (c *DummyClient) Drop(reason string) error {
	log.Println("dropping client", c.Token.Character, ":", reason)
	return nil
}
