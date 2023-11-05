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
	Packet_Which_entityDespawn Packet_Which = 5
)

func (w Packet_Which) String() string {
	const s = "unknownauthentityMoveentitySpawnentityObserveentityDespawn"
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
	case Packet_Which_entityDespawn:
		return s[45:58]

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

func (s Packet) EntityDespawn() (EntityDespawnPacket, error) {
	if capnp.Struct(s).Uint16(0) != 5 {
		panic("Which() != entityDespawn")
	}
	p, err := capnp.Struct(s).Ptr(0)
	return EntityDespawnPacket(p.Struct()), err
}

func (s Packet) HasEntityDespawn() bool {
	if capnp.Struct(s).Uint16(0) != 5 {
		return false
	}
	return capnp.Struct(s).HasPtr(0)
}

func (s Packet) SetEntityDespawn(v EntityDespawnPacket) error {
	capnp.Struct(s).SetUint16(0, 5)
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewEntityDespawn sets the entityDespawn field to a newly
// allocated EntityDespawnPacket struct, preferring placement in s's segment.
func (s Packet) NewEntityDespawn() (EntityDespawnPacket, error) {
	capnp.Struct(s).SetUint16(0, 5)
	ss, err := NewEntityDespawnPacket(capnp.Struct(s).Segment())
	if err != nil {
		return EntityDespawnPacket{}, err
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
func (p Packet_Future) EntityDespawn() EntityDespawnPacket_Future {
	return EntityDespawnPacket_Future{Future: p.Future.Field(0, nil)}
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
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1})
	return EntityMovePacket(st), err
}

func NewRootEntityMovePacket(s *capnp.Segment) (EntityMovePacket, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1})
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

func (s EntityMovePacket) Delta() float32 {
	return math.Float32frombits(capnp.Struct(s).Uint32(16))
}

func (s EntityMovePacket) SetDelta(v float32) {
	capnp.Struct(s).SetUint32(16, math.Float32bits(v))
}

// EntityMovePacket_List is a list of EntityMovePacket.
type EntityMovePacket_List = capnp.StructList[EntityMovePacket]

// NewEntityMovePacket creates a new list of EntityMovePacket.
func NewEntityMovePacket_List(s *capnp.Segment, sz int32) (EntityMovePacket_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1}, sz)
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

type EntityDespawnPacket capnp.Struct

// EntityDespawnPacket_TypeID is the unique identifier for the type EntityDespawnPacket.
const EntityDespawnPacket_TypeID = 0xe597cefc768f55cb

func NewEntityDespawnPacket(s *capnp.Segment) (EntityDespawnPacket, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return EntityDespawnPacket(st), err
}

func NewRootEntityDespawnPacket(s *capnp.Segment) (EntityDespawnPacket, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return EntityDespawnPacket(st), err
}

func ReadRootEntityDespawnPacket(msg *capnp.Message) (EntityDespawnPacket, error) {
	root, err := msg.Root()
	return EntityDespawnPacket(root.Struct()), err
}

func (s EntityDespawnPacket) String() string {
	str, _ := text.Marshal(0xe597cefc768f55cb, capnp.Struct(s))
	return str
}

func (s EntityDespawnPacket) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (EntityDespawnPacket) DecodeFromPtr(p capnp.Ptr) EntityDespawnPacket {
	return EntityDespawnPacket(capnp.Struct{}.DecodeFromPtr(p))
}

func (s EntityDespawnPacket) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s EntityDespawnPacket) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s EntityDespawnPacket) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s EntityDespawnPacket) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s EntityDespawnPacket) Entity() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s EntityDespawnPacket) SetEntity(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

// EntityDespawnPacket_List is a list of EntityDespawnPacket.
type EntityDespawnPacket_List = capnp.StructList[EntityDespawnPacket]

// NewEntityDespawnPacket creates a new list of EntityDespawnPacket.
func NewEntityDespawnPacket_List(s *capnp.Segment, sz int32) (EntityDespawnPacket_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return capnp.StructList[EntityDespawnPacket](l), err
}

// EntityDespawnPacket_Future is a wrapper for a EntityDespawnPacket promised by a client call.
type EntityDespawnPacket_Future struct{ *capnp.Future }

func (f EntityDespawnPacket_Future) Struct() (EntityDespawnPacket, error) {
	p, err := f.Future.Ptr()
	return EntityDespawnPacket(p.Struct()), err
}

const schema_9c4788c5a214e29c = "x\xda\xacSMH\x1cg\x18~\x9f\xef\x9bu\xab8" +
	"\xb8\xe3\xcc\xa1\xd0\xc2\x1e\x0a\xa5\x94ZZ\xbdyY[" +
	"\x94\xfePq?\xac\x95\x96\x1e:\xae\x03\x8a\xed\xcc\xd4" +
	"\x9d]\xb5 V\xda\x82\x05\xa1\xb6XjA\xc1H\x08" +
	"\x0a\x09I \x90\x1frHB\x12\x93H~\x0c\x18\x10" +
	"\x92KH\x0e\x06r\x88Gc\x9c\xf0\xcd\xac;\xbb\xeb" +
	"\x9eBn/\xef\xfb|\x0f\xcf\xfb>\xcf\xf7\xd1}\xb4" +
	")\x1f\xab\xb791\xf1A\xac\xc6?\xf7\xde\xccv[" +
	"fe\x86\x84\x0a\xe6\xcf?4\x96\xaeL}6O1" +
	"\xc4\x89\xf4\x9b\xd8\xd17\x83j\x03\xc7\x09\xfe\x85o\xef" +
	"=i\xf8\xe3\xc6\xaa\x04#\x02w ^C\xa4\xe7\xd8" +
	"\x92>\xce\xe2D-c\xac\x17\x04\xffb\xd3\xda\xfaO" +
	"\xbd[\xab\x15\xdcJ@\xc8\xff\xd1\x1fpYmrI" +
	"\xfd\xd6U\xef\xc8\xdb\xbf\xbds\xad\x82Zb[~V" +
	"\xea\xa0O\x06\xcf\xc6\x95\x14\xc1_\xeb\xf9+\xbf{\xeb" +
	"\xbf\xc7\xd5\xc0\x0b\x12|,\x00\xaf\x04\xe0;\x8fN/" +
	"\xcf\xf6\x8fnK0\xaf\xdc\xf0\xba\xb2\xado\x04\xef\xd6" +
	"\x95\xa4\xd4\x9c\xbf\xbb\xd9\xdd\xda\xf3\xee^\x15j}+" +
	"vI\x7f\x16\x93\xd5\xd3X\x8a\x9a|\xd7\xcc\x0cY\xde" +
	"\x87\x19f\xba\xb6\xdb\xdaa{\x83\xdeX\xb7k\x8e\xd8" +
	"\xa9t0I\x03\xa2\x9e+D\x0a\x88\xb4\x8eV\"\xd1" +
	"\xc6!\xbeb\x00\x0c\xc8\xde\x17_\x12\x89\xcf9\xc4\xd7" +
	"\x0c\x1ac\x06\x18\x91&d3\xcd!\xbegHY\x01" +
	"+j\x89\xa1\x96\xe0\xbbNv\xd0\x1btl\"B\"" +
	":0\x01\x09\x82?\xecx\xe6\xfe\xb4\x8e\x18\xea\xe4\x8b" +
	"\x82H\x04\"\xd3f&^P\xf6&W\xea}?\x90" +
	"\xf6\xff\xa7Db\x96C,2\xa8\xd8\xf3Cq\x0b\xef" +
	"\x13\x899\x0eq\x98Ae/\xfcP\xdd\xa1\xef\x88\xc4" +
	"\"\x878\xca\xa0\xf2]\xdf\x00'\xd2V\xfa\x88\xc42" +
	"\x878\xc5\xa0*\xcf}\x03\x0a\x91vr\x98H\x9c\xe0" +
	"\x10\xe7\x19\xd4\xd8\x8eo F\xa4\x9d\x95\xdd3\x1c\xe2" +
	"2\xc3D\xce\x1e\xb2\x9d\x11\x9bj\x1a\xcc\x9c7\x80D" +
	"d@a\xa7\xf0\x00\x9d\x0e\xf1\xbc\x85D\xe4f\xd9\xb8" +
	"\xdb\xa5\xb89b#\x11\xe5\xb9l\xde\xd5G\xc9\xac5" +
	"\x1c0\x14\x93V\x86h\xb7(\x99uC\x8eb\xbc\x0a" +
	"\x88\xf2\x13~ce\xd0Ramcd\xadV\xf4V" +
	"6\xdb9D\xba\xc4\xdb\xce\xc6\xc8p\x8c\xee\x9b\x84\xb1" +
	"b\xf5\xcb\x01\xe3xI\xba\xba\xfa\x82-\xc2x\x11I" +
	"\x15JQ\x85*\x03\xf6\x06\x870\xaa\xe4\xa6\x0aY\xbb" +
	"\x15,\xfc\xead\xa5\xb9\xeft\xf2V:Y\x8c\xbdQ" +
	"$\x1a\x97D\xa3\x1c\xe2\xf7(\xf6\x932\xe1\xbfr\x88" +
	"\xe9\x92\xd3\xfc)\x9bS\x1cb\x96A\xe3?\x84\xb9\xfa" +
	"[&s\x9aC\xcc1h\x8a\x12\xc6\xea\xdff\"1" +
	"\xc3!\xe6_\xc3\x07\x99\xc8z\x8e\xebZ\xfd\x001\x80" +
	"\x90\xec\xb7~\xf4\xcc\x03.\x84\xbb~\x92\xf3\x06\x0a\xf7" +
	"\xaa8Wst\xae\xa4\xe7\x0cY\xf6\xbe\xa2\x97\x01\x00" +
	"\x00\xff\xff\xcc\x89h\x9f"

func RegisterSchema(reg *schemas.Registry) {
	reg.Register(&schemas.Schema{
		String: schema_9c4788c5a214e29c,
		Nodes: []uint64{
			0x90a96340f29028ba,
			0xc7ca850fead659c0,
			0xc7e9576dd1cb2dc1,
			0xc823831ca674c61b,
			0xe597cefc768f55cb,
			0xf2786494a8b7e4d0,
			0xfe26553a53d9d276,
		},
		Compressed: true,
	})
}
