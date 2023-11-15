package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/net"
	"github.com/johanhenriksson/goworld/game/server"
	"github.com/johanhenriksson/goworld/math/vec3"

	"capnproto.org/go/capnp/v3"
	"github.com/quic-go/quic-go"
)

type Client struct {
	object.Object

	conn   quic.Connection
	stream quic.Stream

	arena   *capnp.SingleSegmentArena
	encoder *capnp.Encoder
	decoder *capnp.Decoder
	submit  func(Event)
}

func NewClient(events func(Event)) *Client {
	client := object.New("GameClient", &Client{
		arena:  capnp.SingleSegment(nil),
		submit: events,
	})
	return client
}

func (c *Client) Connect(hostname string) error {
	if c.conn != nil {
		return fmt.Errorf("already connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := quic.DialAddr(
		ctx,
		fmt.Sprintf("%s:%d", hostname, net.GamePort),
		&tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"goworld"},
		},
		&quic.Config{
			KeepAlivePeriod: 3 * time.Second,
			EnableDatagrams: true,
		},
	)
	if err != nil {
		return err
	}

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return err
	}

	c.conn = conn
	c.stream = stream
	c.encoder = capnp.NewEncoder(stream)
	c.decoder = capnp.NewDecoder(stream)
	go c.readLoop()

	return nil
}

func (c *Client) Disconnect() error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	return c.conn.CloseWithError(0, "disconnected")
}

func (c *Client) decode() (*net.Packet, error) {
	msg, err := c.decoder.Decode()
	if err != nil {
		return nil, err
	}

	pkt, err := net.ReadRootPacket(msg)
	if err != nil {
		return nil, err
	}

	return &pkt, nil
}

func (c *Client) handlePacket(msg *net.Packet) error {
	defer msg.Message().Release()

	switch msg.Which() {
	case net.Packet_Which_auth:
		return fmt.Errorf("%w: auth packet is client->server only", net.ErrInvalidPacket)

	case net.Packet_Which_enterWorld:
		enter, err := msg.EnterWorld()
		if err != nil {
			return err
		}

		mapName, err := enter.Map()
		if err != nil {
			return err
		}

		c.submit(EnterWorldEvent{
			Map: mapName,
		})

	case net.Packet_Which_entityMove:
		move, err := msg.EntityMove()
		if err != nil {
			return err
		}
		pos, err := move.Position()
		if err != nil {
			return err
		}

		c.submit(EntityMoveEvent{
			EntityID: server.Identity(move.Entity()),
			Position: net.ToVec3(pos),
			Rotation: move.Rotation(),
			Stopped:  move.Stopped(),
			Delta:    move.Delta(),
		})

	case net.Packet_Which_entitySpawn:
		spawn, err := msg.EntitySpawn()
		if err != nil {
			return err
		}
		pos, err := spawn.Position()
		if err != nil {
			return err
		}
		c.submit(EntitySpawnEvent{
			EntityID: server.Identity(spawn.Entity()),
			Position: net.ToVec3(pos),
			Rotation: spawn.Rotation(),
		})

	case net.Packet_Which_entityDespawn:
		despawn, err := msg.EntityDespawn()
		if err != nil {
			return err
		}
		c.submit(EntityDespawnEvent{
			EntityID: server.Identity(despawn.Entity()),
		})

	case net.Packet_Which_entityObserve:
		observe, err := msg.EntityObserve()
		if err != nil {
			return err
		}
		c.submit(EntityObserveEvent{
			EntityID: server.Identity(observe.Entity()),
		})

	default:
		return fmt.Errorf("%w: received packet with type %s", net.ErrUnknownPacket, msg.Which())
	}

	return nil
}

func (c *Client) readLoop() {
	defer func() {
		// disconnected. reset client:
		c.conn = nil
		c.encoder = nil
		c.decoder = nil

		c.submit(DisconnectEvent{})

		log.Println("disconnected from server")
	}()

	for {
		msg, err := c.decode()
		if err != nil {
			log.Println("client read error:", err)
			return
		}

		if err := c.handlePacket(msg); err != nil {
			log.Println("client packet error:", err)
			return
		}
	}
}

type PacketBuilderFn func(*net.Packet) error

func (c *Client) Send(fn PacketBuilderFn) error {
	// allocate message
	msg, seg, err := capnp.NewMessage(c.arena)
	if err != nil {
		return err
	}

	// remember to clean up after transmission
	defer msg.Release()

	// create packet wrapper
	wrap, err := net.NewRootPacket(seg)
	if err != nil {
		return err
	}

	// delegate to packet builder func
	if err := fn(&wrap); err != nil {
		return err
	}

	return c.encoder.Encode(msg)
}

func (c *Client) SendAuthToken(token uint64) error {
	return c.Send(func(p *net.Packet) error {
		auth, err := net.NewAuthPacket(p.Segment())
		if err != nil {
			return err
		}

		// encode authentication token
		auth.SetToken(token)

		return p.SetAuth(auth)
	})
}

func (c *Client) SendMove(id server.Identity, position vec3.T, rotation float32, stopped bool, delta float32) error {
	return c.Send(func(p *net.Packet) error {
		move, err := net.NewEntityMovePacket(p.Segment())
		if err != nil {
			return err
		}

		move.SetEntity(uint64(id))

		pos, err := net.FromVec3(p.Segment(), position)
		if err != nil {
			return err
		}
		move.SetPosition(pos)

		move.SetRotation(rotation)

		move.SetStopped(stopped)

		move.SetDelta(delta)

		return p.SetEntityMove(move)
	})
}
