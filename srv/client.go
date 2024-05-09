package srv

import (
	"log"
)

type ClientToken struct {
	Character string
}

type Client interface {
	Observe(Area, Actor)

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
	Area  Area
	Actor Actor
}

var _ Client = (*DummyClient)(nil)

func (c *DummyClient) ReadToken() (ClientToken, error) {
	return c.Token, nil
}

func (c *DummyClient) Send(ev Event) {
	// do nothing
	log.Println(c.Token.Character, "<-", Dump(ev))
}

func (c *DummyClient) Drop(reason string) error {
	if c.Area != nil {
		c.Area.Unsubscribe(c)
		c.Area = nil
		c.Actor = nil
	}
	log.Println("dropping client", c.Token.Character, ":", reason)
	return nil
}

func (c *DummyClient) Observe(area Area, actor Actor) {
	c.Area = area
	c.Actor = actor
	log.Println("client", c.Token.Character, "observing", actor.Name())
	c.Area.Subscribe(c, func(ev Event) {
		c.Send(ev)
	})
}

func (c *DummyClient) Action(a Action) {
	c.Actor.Area().Action(a)
}
