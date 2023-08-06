package uniform

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

type Light struct {
	ViewProj    [4]mat4.T
	Shadowmap   [4]uint32
	Distance    [4]float32
	Color       color.T
	Position    vec4.T
	Type        uint32
	Intensity   float32
	Range       float32
	Attenuation Attenuation
}

type Attenuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

type LightSettings struct {
	AmbientColor       color.T
	AmbientIntensity   float32
	Count              int32
	ShadowSamples      int32
	ShadowSampleRadius float32
	ShadowBias         float32
	NormalOffset       float32
	_padding           [76]float32
}

func init() {
	lightSz := unsafe.Sizeof(Light{})
	settingSz := unsafe.Sizeof(LightSettings{})
	if lightSz != settingSz {
		panic(fmt.Sprintf("Light (%d) and LightSetting (%d) must have equal size", lightSz, settingSz))
	}
}
