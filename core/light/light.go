package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

type ShadowmapStore interface {
	Lookup(T, int) (int, bool)
}

type T interface {
	object.Component

	Type() Type
	CastShadows() bool
	Shadowmaps() int
	LightData(ShadowmapStore) uniform.Light
	ShadowProjection(mapIndex int) uniform.Camera
}

// Descriptor holds rendering information for lights
type Descriptor struct {
	Projection mat4.T // Light projection matrix
	View       mat4.T // Light view matrix
	ViewProj   mat4.T
	Color      color.T
	Position   vec4.T
	Type       Type
	Range      float32
	Intensity  float32
	Shadows    uint32
	Index      int
}
