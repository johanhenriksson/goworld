package server

import "github.com/johanhenriksson/goworld/math/vec3"

type Player struct {
	instance *Instance
}

var _ Entity = &Player{}

func (p *Player) ID() Identity {
	return 0
}

func (p *Player) Name() string {
	return "Player 1"
}

func (p *Player) Instance() *Instance {
	return p.instance
}

func (p *Player) Position() vec3.T {
	return vec3.Zero
}
