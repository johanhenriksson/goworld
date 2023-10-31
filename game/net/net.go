package net

import capnp "capnproto.org/go/capnp/v3"

type Packet interface {
	Message() *capnp.Message
}
