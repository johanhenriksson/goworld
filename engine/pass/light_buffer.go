package pass

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type ShadowmapLookupFn func(light.T, int) *texture.Texture

type LightBuffer struct {
	buffer   []uniform.Light
	settings uniform.LightSettings
}

func NewLightBuffer(capacity int) *LightBuffer {
	return &LightBuffer{
		buffer: make([]uniform.Light, 1, capacity+1),

		// default lighting settings
		settings: uniform.LightSettings{
			AmbientColor:     color.White,
			AmbientIntensity: 0.4,

			ShadowBias:         0.005,
			ShadowSampleRadius: 1,
			ShadowSamples:      1,
			NormalOffset:       0.1,
		},
	}
}

func (b *LightBuffer) Flush(desc *descriptor.Storage[uniform.Light]) {
	// settings is stored in the first element of the buffer
	// it excludes the first element containing the light settings
	b.settings.Count = int32(len(b.buffer) - 1)
	b.buffer[0] = *(*uniform.Light)(unsafe.Pointer(&b.settings))
	desc.SetRange(0, b.buffer)
}

func (b *LightBuffer) Reset() {
	b.buffer = b.buffer[:1]
}

func (b *LightBuffer) Store(light uniform.Light) {
	b.buffer = append(b.buffer, light)
}
