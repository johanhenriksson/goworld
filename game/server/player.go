package server

import "github.com/johanhenriksson/goworld/math/vec3"

type Player struct {
	id Identity
}

var _ Entity = &Player{}

func (p *Player) ID() Identity {
	return p.id
}

func (p *Player) Name() string {
	return "Player 1"
}

func (p *Player) Position() vec3.T {
	return vec3.New(1, 2, 3)
}
