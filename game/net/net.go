package net

import (
	"errors"

	capnp "capnproto.org/go/capnp/v3"
	"github.com/johanhenriksson/goworld/math/vec3"
)

const GamePort = 1423

var ErrUnknownPacket = errors.New("unknown packet")
var ErrInvalidPacket = errors.New("invalid packet")

type Message interface {
	Message() *capnp.Message
}

func FromVec3(seg *capnp.Segment, v vec3.T) (Vec3, error) {
	vo, err := NewVec3(seg)
	if err != nil {
		return Vec3{}, err
	}
	vo.SetX(v.X)
	vo.SetY(v.Y)
	vo.SetZ(v.Z)
	return vo, nil
}

func ToVec3(v Vec3) vec3.T {
	return vec3.New(v.X(), v.Y(), v.Z())
}
