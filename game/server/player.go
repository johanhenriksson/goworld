package server

import "github.com/johanhenriksson/goworld/math/vec3"

type Player struct {
	id       Identity
	position vec3.T
	rotation float32
}

var _ Entity = &Player{}

func (p *Player) ID() Identity {
	return p.id
}

func (p *Player) Name() string {
	return "Player 1"
}

func (p *Player) Position() vec3.T       { return p.position }
func (p *Player) SetPosition(pos vec3.T) { p.position = pos }

func (p *Player) Rotation() float32       { return p.rotation }
func (p *Player) SetRotation(rot float32) { p.rotation = rot }
