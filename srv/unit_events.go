package srv

import "github.com/johanhenriksson/goworld/math/vec3"

type PositionUpdateEvent struct {
	Unit     Actor
	Position vec3.T
}

func (e PositionUpdateEvent) Source() any { return e.Unit }

type TargetUpdateEvent struct {
	Unit   Actor
	Target Identity
}

func (e TargetUpdateEvent) Source() any { return e.Unit }
