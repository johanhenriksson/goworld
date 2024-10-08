package uniform

import (
	"fmt"
	"structs"
	"unsafe"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

const ShadowCascades = 4
const LightPadding = 74

type Light struct {
	_ structs.HostLayout

	ViewProj  [ShadowCascades]mat4.T
	Shadowmap [ShadowCascades]uint32
	Distance  [ShadowCascades]float32
	Color     color.T
	Position  vec4.T
	Type      uint32
	Intensity float32
	Range     float32
	Falloff   float32
}

type LightSettings struct {
	_ structs.HostLayout

	AmbientColor       color.T
	AmbientIntensity   float32
	Count              int32
	ShadowSamples      int32
	ShadowSampleRadius float32
	ShadowBias         float32
	NormalOffset       float32
	_padding           [LightPadding]uint32
}

func init() {
	lightSz := unsafe.Sizeof(Light{})
	settingSz := unsafe.Sizeof(LightSettings{})
	if lightSz != settingSz {
		panic(fmt.Sprintf("Light (%d) and LightSetting (%d) must have equal size", lightSz, settingSz))
	}
}
