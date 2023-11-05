// Code generated by capnpc-go. DO NOT EDIT.

package net

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
	math "math"
	strconv "strconv"
)

type Vec3 capnp.Struct

// Vec3_TypeID is the unique identifier for the type Vec3.
const Vec3_TypeID = 0xc7e9576dd1cb2dc1

func NewVec3(s *capnp.Segment) (Vec3, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 0})
	return Vec3(st), err
}

func NewRootVec3(s *capnp.Segment) (Vec3, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 0})
	return Vec3(st), err
}

func ReadRootVec3(msg *capnp.Message) (Vec3, error) {
	root, err := msg.Root()
	return Vec3(root.Struct()), err
}

func (s Vec3) String() string {
	str, _ := text.Marshal(0xc7e9576dd1cb2dc1, capnp.Struct(s))
	return str
}

func (s Vec3) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Vec3) DecodeFromPtr(p capnp.Ptr) Vec3 {
	return Vec3(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Vec3) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Vec3) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Vec3) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Vec3) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Vec3) X() float32 {
	return math.Float32frombits(capnp.Struct(s).Uint32(0))
}

func (s Vec3) SetX(v float32) {
	capnp.Struct(s).SetUint32(0, math.Float32bits(v))
}

func (s Vec3) Y() float32 {
	return math.Float32frombits(capnp.Struct(s).Uint32(4))
}

func (s Vec3) SetY(v float32) {
	capnp.Struct(s).SetUint32(4, math.Float32bits(v))
}

func (s Vec3) Z() float32 {
	return math.Float32frombits(capnp.Struct(s).Uint32(8))
}

func (s Vec3) SetZ(v float32) {
	capnp.Struct(s).SetUint32(8, math.Float32bits(v))
}

// Vec3_List is a list of Vec3.
type Vec3_List = capnp.StructList[Vec3]

// NewVec3 creates a new list of Vec3.
func NewVec3_List(s *capnp.Segment, sz int32) (Vec3_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 0}, sz)
	return capnp.StructList[Vec3](l), err
}

// Vec3_Future is a wrapper for a Vec3 promised by a client call.
type Vec3_Future struct{ *capnp.Future }

func (f Vec3_Future) Struct() (Vec3, error) {
	p, err := f.Future.Ptr()
	return Vec3(p.Struct()), err
}

type Packet capnp.Struct
type Packet_Which uint16

const (
	Packet_Which_unknown       Packet_Which = 0
	Packet_Which_auth          Packet_Which = 1
	Packet_Which_entityMove    Packet_Which = 2
	Packet_Which_entitySpawn   Packet_Which = 3
	Packet_Which_entityObserve Packet_Which = 4
)

func (w Packet_Which) String() string {
	const s = "unknownauthentityMoveentitySpawnentityObserve"
	switch w {
	case Packet_Which_unknown:
		return s[0:7]
	case Packet_Which_auth:
		return s[7:11]
	case Packet_Which_entityMove:
		return s[11:21]
	case Packet_Which_entitySpawn:
		return s[21:32]
	case Packet_Which_entityObserve:
		return s[32:45]

	}
	return "Packet_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Packet_TypeID is the unique identifier for the type Packet.
const Packet_TypeID = 0xc7ca850fead659c0

func NewPacket(s *capnp.Segment) (Packet, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Packet(st), err
}

func NewRootPacket(s *capnp.Segment) (Packet, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Packet(st), err
}

func ReadRootPacket(msg *capnp.Message) (Packet, error) {
	root, err := msg.Root()
	return Packet(root.Struct()), err
}

func (s Packet) String() string {
	str, _ := text.Marshal(0xc7ca850fead659c0, capnp.Struct(s))
	return str
}

func (s Packet) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Packet) DecodeFromPtr(p capnp.Ptr) Packet {
	return Packet(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Packet) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}

func (s Packet) Which() Packet_Which {
	return Packet_Which(capnp.Struct(s).Uint16(0))
}
func (s Packet) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Packet) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Packet) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Packet) SetUnknown() {
	capnp.Struct(s).SetUint16(0, 0)

}

func (s Packet) Auth() (AuthPacket, error) {
	if capnp.Struct(s).Uint16(0) != 1 {
		panic("Which() != auth")
	}
	p, err := capnp.Struct(s).Ptr(0)
	return AuthPacket(p.Struct()), err
}

func (s Packet) HasAuth() bool {
	if capnp.Struct(s).Uint16(0) != 1 {
		return false
	}
	return capnp.Struct(s).HasPtr(0)
}

func (s Packet) SetAuth(v AuthPacket) error {
	capnp.Struct(s).SetUint16(0, 1)
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewAuth sets the auth field to a newly
// allocated AuthPacket struct, preferring placement in s's segment.
func (s Packet) NewAuth() (AuthPacket, error) {
	capnp.Struct(s).SetUint16(0, 1)
	ss, err := NewAuthPacket(capnp.Struct(s).Segment())
	if err != nil {
		return AuthPacket{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Packet) EntityMove() (EntityMovePacket, error) {
	if capnp.Struct(s).Uint16(0) != 2 {
		panic("Which() != entityMove")
	}
	p, err := capnp.Struct(s).Ptr(0)
	return EntityMovePacket(p.Struct()), err
}

func (s Packet) HasEntityMove() bool {
	if capnp.Struct(s).Uint16(0) != 2 {
		return false
	}
	return capnp.Struct(s).HasPtr(0)
}

func (s Packet) SetEntityMove(v EntityMovePacket) error {
	capnp.Struct(s).SetUint16(0, 2)
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewEntityMove sets the entityMove field to a newly
// allocated EntityMovePacket struct, preferring placement in s's segment.
func (s Packet) NewEntityMove() (EntityMovePacket, error) {
	capnp.Struct(s).SetUint16(0, 2)
	ss, err := NewEntityMovePacket(capnp.Struct(s).Segment())
	if err != nil {
		return EntityMovePacket{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Packet) EntitySpawn() (EntitySpawnPacket, error) {
	if capnp.Struct(s).Uint16(0) != 3 {
		panic("Which() != entitySpawn")
	}
	p, err := capnp.Struct(s).Ptr(0)
	return EntitySpawnPacket(p.Struct()), err
}

func (s Packet) HasEntitySpawn() bool {
	if capnp.Struct(s).Uint16(0) != 3 {
		return false
	}
	return capnp.Struct(s).HasPtr(0)
}

func (s Packet) SetEntitySpawn(v EntitySpawnPacket) error {
	capnp.Struct(s).SetUint16(0, 3)
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewEntitySpawn sets the entitySpawn field to a newly
// allocated EntitySpawnPacket struct, preferring placement in s's segment.
func (s Packet) NewEntitySpawn() (EntitySpawnPacket, error) {
	capnp.Struct(s).SetUint16(0, 3)
	ss, err := NewEntitySpawnPacket(capnp.Struct(s).Segment())
	if err != nil {
		return EntitySpawnPacket{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Packet) EntityObserve() (EntityObservePacket, error) {
	if capnp.Struct(s).Uint16(0) != 4 {
		panic("Which() != entityObserve")
	}
	p, err := capnp.Struct(s).Ptr(0)
	return EntityObservePacket(p.Struct()), err
}

func (s Packet) HasEntityObserve() bool {
	if capnp.Struct(s).Uint16(0) != 4 {
		return false
	}
	return capnp.Struct(s).HasPtr(0)
}

func (s Packet) SetEntityObserve(v EntityObservePacket) error {
	capnp.Struct(s).SetUint16(0, 4)
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewEntityObserve sets the entityObserve field to a newly
// allocated EntityObservePacket struct, preferring placement in s's segment.
func (s Packet) NewEntityObserve() (EntityObservePacket, error) {
	capnp.Struct(s).SetUint16(0, 4)
	ss, err := NewEntityObservePacket(capnp.Struct(s).Segment())
	if err != nil {
		return EntityObservePacket{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

// Packet_List is a list of Packet.
type Packet_List = capnp.StructList[Packet]

// NewPacket creates a new list of Packet.
func NewPacket_List(s *capnp.Segment, sz int32) (Packet_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return capnp.StructList[Packet](l), err
}

// Packet_Future is a wrapper for a Packet promised by a client call.
type Packet_Future struct{ *capnp.Future }

func (f Packet_Future) Struct() (Packet, error) {
	p, err := f.Future.Ptr()
	return Packet(p.Struct()), err
}
func (p Packet_Future) Auth() AuthPacket_Future {
	return AuthPacket_Future{Future: p.Future.Field(0, nil)}
}
func (p Packet_Future) EntityMove() EntityMovePacket_Future {
	return EntityMovePacket_Future{Future: p.Future.Field(0, nil)}
}
func (p Packet_Future) EntitySpawn() EntitySpawnPacket_Future {
	return EntitySpawnPacket_Future{Future: p.Future.Field(0, nil)}
}
func (p Packet_Future) EntityObserve() EntityObservePacket_Future {
	return EntityObservePacket_Future{Future: p.Future.Field(0, nil)}
}

type AuthPacket capnp.Struct

// AuthPacket_TypeID is the unique identifier for the type AuthPacket.
const AuthPacket_TypeID = 0xfe26553a53d9d276

func NewAuthPacket(s *capnp.Segment) (AuthPacket, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return AuthPacket(st), err
}

func NewRootAuthPacket(s *capnp.Segment) (AuthPacket, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return AuthPacket(st), err
}

func ReadRootAuthPacket(msg *capnp.Message) (AuthPacket, error) {
	root, err := msg.Root()
	return AuthPacket(root.Struct()), err
}

func (s AuthPacket) String() string {
	str, _ := text.Marshal(0xfe26553a53d9d276, capnp.Struct(s))
	return str
}

func (s AuthPacket) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (AuthPacket) DecodeFromPtr(p capnp.Ptr) AuthPacket {
	return AuthPacket(capnp.Struct{}.DecodeFromPtr(p))
}

func (s AuthPacket) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s AuthPacket) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s AuthPacket) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s AuthPacket) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s AuthPacket) Token() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s AuthPacket) SetToken(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

// AuthPacket_List is a list of AuthPacket.
type AuthPacket_List = capnp.StructList[AuthPacket]

// NewAuthPacket creates a new list of AuthPacket.
func NewAuthPacket_List(s *capnp.Segment, sz int32) (AuthPacket_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return capnp.StructList[AuthPacket](l), err
}

// AuthPacket_Future is a wrapper for a AuthPacket promised by a client call.
type AuthPacket_Future struct{ *capnp.Future }

func (f AuthPacket_Future) Struct() (AuthPacket, error) {
	p, err := f.Future.Ptr()
	return AuthPacket(p.Struct()), err
}

type EntityMovePacket capnp.Struct

// EntityMovePacket_TypeID is the unique identifier for the type EntityMovePacket.
const EntityMovePacket_TypeID = 0xf2786494a8b7e4d0

func NewEntityMovePacket(s *capnp.Segment) (EntityMovePacket, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return EntityMovePacket(st), err
}

func NewRootEntityMovePacket(s *capnp.Segment) (EntityMovePacket, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return EntityMovePacket(st), err
}

func ReadRootEntityMovePacket(msg *capnp.Message) (EntityMovePacket, error) {
	root, err := msg.Root()
	return EntityMovePacket(root.Struct()), err
}

func (s EntityMovePacket) String() string {
	str, _ := text.Marshal(0xf2786494a8b7e4d0, capnp.Struct(s))
	return str
}

func (s EntityMovePacket) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (EntityMovePacket) DecodeFromPtr(p capnp.Ptr) EntityMovePacket {
	return EntityMovePacket(capnp.Struct{}.DecodeFromPtr(p))
}

func (s EntityMovePacket) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s EntityMovePacket) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s EntityMovePacket) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s EntityMovePacket) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s EntityMovePacket) Entity() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s EntityMovePacket) SetEntity(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s EntityMovePacket) Position() (Vec3, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return Vec3(p.Struct()), err
}

func (s EntityMovePacket) HasPosition() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s EntityMovePacket) SetPosition(v Vec3) error {
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewPosition sets the position field to a newly
// allocated Vec3 struct, preferring placement in s's segment.
func (s EntityMovePacket) NewPosition() (Vec3, error) {
	ss, err := NewVec3(capnp.Struct(s).Segment())
	if err != nil {
		return Vec3{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s EntityMovePacket) Rotation() float32 {
	return math.Float32frombits(capnp.Struct(s).Uint32(8))
}

func (s EntityMovePacket) SetRotation(v float32) {
	capnp.Struct(s).SetUint32(8, math.Float32bits(v))
}

func (s EntityMovePacket) Stopped() bool {
	return capnp.Struct(s).Bit(96)
}

func (s EntityMovePacket) SetStopped(v bool) {
	capnp.Struct(s).SetBit(96, v)
}

// EntityMovePacket_List is a list of EntityMovePacket.
type EntityMovePacket_List = capnp.StructList[EntityMovePacket]

// NewEntityMovePacket creates a new list of EntityMovePacket.
func NewEntityMovePacket_List(s *capnp.Segment, sz int32) (EntityMovePacket_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	return capnp.StructList[EntityMovePacket](l), err
}

// EntityMovePacket_Future is a wrapper for a EntityMovePacket promised by a client call.
type EntityMovePacket_Future struct{ *capnp.Future }

func (f EntityMovePacket_Future) Struct() (EntityMovePacket, error) {
	p, err := f.Future.Ptr()
	return EntityMovePacket(p.Struct()), err
}
func (p EntityMovePacket_Future) Position() Vec3_Future {
	return Vec3_Future{Future: p.Future.Field(0, nil)}
}

type EntitySpawnPacket capnp.Struct

// EntitySpawnPacket_TypeID is the unique identifier for the type EntitySpawnPacket.
const EntitySpawnPacket_TypeID = 0x90a96340f29028ba

func NewEntitySpawnPacket(s *capnp.Segment) (EntitySpawnPacket, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return EntitySpawnPacket(st), err
}

func NewRootEntitySpawnPacket(s *capnp.Segment) (EntitySpawnPacket, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return EntitySpawnPacket(st), err
}

func ReadRootEntitySpawnPacket(msg *capnp.Message) (EntitySpawnPacket, error) {
	root, err := msg.Root()
	return EntitySpawnPacket(root.Struct()), err
}

func (s EntitySpawnPacket) String() string {
	str, _ := text.Marshal(0x90a96340f29028ba, capnp.Struct(s))
	return str
}

func (s EntitySpawnPacket) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (EntitySpawnPacket) DecodeFromPtr(p capnp.Ptr) EntitySpawnPacket {
	return EntitySpawnPacket(capnp.Struct{}.DecodeFromPtr(p))
}

func (s EntitySpawnPacket) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s EntitySpawnPacket) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s EntitySpawnPacket) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s EntitySpawnPacket) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s EntitySpawnPacket) Entity() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s EntitySpawnPacket) SetEntity(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s EntitySpawnPacket) Position() (Vec3, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return Vec3(p.Struct()), err
}

func (s EntitySpawnPacket) HasPosition() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s EntitySpawnPacket) SetPosition(v Vec3) error {
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewPosition sets the position field to a newly
// allocated Vec3 struct, preferring placement in s's segment.
func (s EntitySpawnPacket) NewPosition() (Vec3, error) {
	ss, err := NewVec3(capnp.Struct(s).Segment())
	if err != nil {
		return Vec3{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s EntitySpawnPacket) Rotation() float32 {
	return math.Float32frombits(capnp.Struct(s).Uint32(8))
}

func (s EntitySpawnPacket) SetRotation(v float32) {
	capnp.Struct(s).SetUint32(8, math.Float32bits(v))
}

// EntitySpawnPacket_List is a list of EntitySpawnPacket.
type EntitySpawnPacket_List = capnp.StructList[EntitySpawnPacket]

// NewEntitySpawnPacket creates a new list of EntitySpawnPacket.
func NewEntitySpawnPacket_List(s *capnp.Segment, sz int32) (EntitySpawnPacket_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	return capnp.StructList[EntitySpawnPacket](l), err
}

// EntitySpawnPacket_Future is a wrapper for a EntitySpawnPacket promised by a client call.
type EntitySpawnPacket_Future struct{ *capnp.Future }

func (f EntitySpawnPacket_Future) Struct() (EntitySpawnPacket, error) {
	p, err := f.Future.Ptr()
	return EntitySpawnPacket(p.Struct()), err
}
func (p EntitySpawnPacket_Future) Position() Vec3_Future {
	return Vec3_Future{Future: p.Future.Field(0, nil)}
}

type EntityObservePacket capnp.Struct

// EntityObservePacket_TypeID is the unique identifier for the type EntityObservePacket.
const EntityObservePacket_TypeID = 0xc823831ca674c61b

func NewEntityObservePacket(s *capnp.Segment) (EntityObservePacket, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return EntityObservePacket(st), err
}

func NewRootEntityObservePacket(s *capnp.Segment) (EntityObservePacket, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return EntityObservePacket(st), err
}

func ReadRootEntityObservePacket(msg *capnp.Message) (EntityObservePacket, error) {
	root, err := msg.Root()
	return EntityObservePacket(root.Struct()), err
}

func (s EntityObservePacket) String() string {
	str, _ := text.Marshal(0xc823831ca674c61b, capnp.Struct(s))
	return str
}

func (s EntityObservePacket) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (EntityObservePacket) DecodeFromPtr(p capnp.Ptr) EntityObservePacket {
	return EntityObservePacket(capnp.Struct{}.DecodeFromPtr(p))
}

func (s EntityObservePacket) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s EntityObservePacket) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s EntityObservePacket) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s EntityObservePacket) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s EntityObservePacket) Entity() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s EntityObservePacket) SetEntity(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

// EntityObservePacket_List is a list of EntityObservePacket.
type EntityObservePacket_List = capnp.StructList[EntityObservePacket]

// NewEntityObservePacket creates a new list of EntityObservePacket.
func NewEntityObservePacket_List(s *capnp.Segment, sz int32) (EntityObservePacket_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return capnp.StructList[EntityObservePacket](l), err
}

// EntityObservePacket_Future is a wrapper for a EntityObservePacket promised by a client call.
type EntityObservePacket_Future struct{ *capnp.Future }

func (f EntityObservePacket_Future) Struct() (EntityObservePacket, error) {
	p, err := f.Future.Ptr()
	return EntityObservePacket(p.Struct()), err
}

const schema_9c4788c5a214e29c = "x\xda\xac\x92OH\x14a\x18\xc6\xdf\xe7\xfbf\xdc\x14" +
	"\x07w\x9c\xf1\x10\x15{\x08B\x02#\xf3\x94\x10k\x91" +
	"\x94\xd1\xe2~\xd8\x1f\x0a\x0f\x8d\xeb\x80\x8b93\xe9\xac" +
	"\xbb\x86!\x91\x82\x87\x02\x09\x83\x0e\x1e\xcaC\xe4!:" +
	"\x06\x9d*(\xa3.\xd5!\xe8\x1aA\xd4\xd1c\x7f\xfc" +
	"\xe2\x9b]w\xdcu\x8f\xdd^\xde\xf7\xfd\x1e~\xdf\xfb" +
	"<\x87W\xd0\xa7u\x1bENLt\xeaM\xf2y\xe7" +
	"\xd2F_nm\x89\x84\x01&W\xbe\xda\xab\xaf\x17O" +
	"\xad\x90\x8e\x04\x91\xf5\x10\xbf\xac'Q\xb5\x86\xa7\x04\xf9" +
	"\xe2\xd2\xe7\x9fm\x0b\xef\xd6\xd52\xe2\xe5~$t\"" +
	"k\x80\xadZ\x82%\x88z2,\x05\x82|\xd9\xf5\xfe" +
	"\xe3\xc4\xc5\x1f\xebu\xda\x9a\x12\xcc\xf3\xbb\xd65\xae\xaa" +
	"\x09\xae\xa4\xf7\xbc\x09\x1f\xed\xbd\xb5\xffm\x9d\xb4\xda\xed" +
	"\xd1\xb5\x16X\x1d\xd13SK\x13\xe4\x87o\xcf\x1e/" +
	"\x8f\x966\x1aBwk\x1b\xd6\xb1h\xf9\xa8\xf6\x9d " +
	"\xa7?}\x19\xea=\x7f`\xb3\x81\xb2\xd5\xa1\xbf\xb2\xf6" +
	"\xe9\xaa\xda\xad\xa7\xa9K\x06Nn\xdc\x0d\x0f\xe5\x98\x13" +
	"xAo\xbf\x17\xe6\xc3\x99\xa1\xc0)z\xe9l4\xc9" +
	"\x02\xa2\x95kD\x1a\x88\xcc\xfe^\"\xd1\xc7!\xce2" +
	"\x006To\xe0\x0c\x918\xcd!\xce1\x98\x8c\xd9`" +
	"D\xa6P\xcd,\x87\x18fH\xbb\x91*\x9a\x89\xa1\x99" +
	" \x03\x7f*\x1f\xe6}\x8f\x88\x90\x8cOF@\x92 " +
	"'\xfd\xd0\xd9\x9a\xb6\x10C\x8bzQ\x81D\x04\x99u" +
	"r\x89\x0a\x99\xcd\xb5V)#\xb4\x1b'\x88D\x89C" +
	"\xcc3\x18\xd8\x94e\xb8\x9b\x07\x89\xc4,\x87Xd0" +
	"\xd8_Y\xa6[\xb8L$\xe69\xc4\x12\x83\xc1\xffH" +
	"\x1b\x9c\xc8\xbc3B$ns\x88\xfb\x0c\x86\xf6[\xda" +
	"\xd0\x88\xcc{\x93Db\x99C<`\x98+x\xe3\x9e" +
	"_\xf4\xa8\xa9\xcd)\x84cH\xc6\xa7\xae\xd0\x97\xbf\x9a" +
	"\xf1\x89O\xbbH\xc6\xb6\xd5\x8c\x87\x02J8E\x0f\xc9" +
	"8\x8b5\xf3\xc1\x11JM\xb9\x93\x91B5%\x95\x8d" +
	"\xdaS\\ps\xe8\xa9\xb3\xa8=\xb6\xc8\xacz\xa4\x9a" +
	"'9Dv\x9bG\x99\xf6\xd88\x94\xb6\x8e\x8d\x99j" +
	"u}\x87\x01|[J\x06G\"\xc6rL\x88\x14\x85" +
	"V\xa50TPvq\x08\xbb\x81\xff\x0d\"\x97\xf1\xa7" +
	"\xddl\xaa\x9a\xb8dU\xc8QB\xc3\x1cb,N\x9c" +
	"\xab\xc25\xca!\x82m\xbf\x99P\xcd\xab\x1c\xa2\xc4`" +
	"\xf2+eK\x0b*\x14\x01\x87\x98\xfd\x0f1\x9c\x9b\x0a" +
	"\xfd pG\x01b\xc0\x8e\x8f\x1c/\x84c\x95c\xd4" +
	"\xdd\xe2H|\x8bT\xe8\x8f\xbb\xde\x16\xc3\xbf\x00\x00\x00" +
	"\xff\xff\xe3'$\x9e"

func RegisterSchema(reg *schemas.Registry) {
	reg.Register(&schemas.Schema{
		String: schema_9c4788c5a214e29c,
		Nodes: []uint64{
			0x90a96340f29028ba,
			0xc7ca850fead659c0,
			0xc7e9576dd1cb2dc1,
			0xc823831ca674c61b,
			0xf2786494a8b7e4d0,
			0xfe26553a53d9d276,
		},
		Compressed: true,
	})
}
