package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/texture"
)

type ShadowmapLookupFn func(light.T, int) texture.T

type LightBuffer struct {
	lights       *descriptor.Storage[uniform.Light]
	shadows      cache.SamplerCache
	buffer       []uniform.Light
	lookupShadow ShadowmapLookupFn
	next         int
}

func NewLightBuffer(lights *descriptor.Storage[uniform.Light], shadowCache cache.SamplerCache, lookup ShadowmapLookupFn) *LightBuffer {
	return &LightBuffer{
		lights:       lights,
		shadows:      shadowCache,
		buffer:       make([]uniform.Light, lights.Size),
		lookupShadow: lookup,
		next:         0,
	}
}

func (b *LightBuffer) Flush() {
	b.lights.SetRange(0, b.buffer[:b.next])
	b.shadows.UpdateDescriptors()
}

func (b *LightBuffer) Reset() {
	b.next = 0
}

func (b *LightBuffer) Count() int {
	return b.next
}

func (b *LightBuffer) Store(args render.Args, lit light.T) {
	desc := lit.LightDescriptor(args, 0)

	entry := uniform.Light{
		Type:      desc.Type,
		Color:     desc.Color,
		Position:  desc.Position,
		Intensity: desc.Intensity,
	}

	switch lit.(type) {
	case *light.Point:
		entry.Attenuation = desc.Attenuation
		entry.Range = desc.Range

	case *light.Directional:
		for cascadeIndex, cascade := range lit.Cascades() {
			entry.ViewProj[cascadeIndex] = cascade.ViewProj
			entry.Distance[cascadeIndex] = cascade.FarSplit

			if shadowtex := b.lookupShadow(lit, cascadeIndex); shadowtex != nil {
				handle := b.shadows.Assign(shadowtex)
				entry.Shadowmap[cascadeIndex] = uint32(handle.ID)
			} else {
				// no shadowmap available - disable shadows until its available
				log.Println("missing cascade shadowmap", cascadeIndex)
				entry.Shadowmap[cascadeIndex] = 0
			}
		}
	}

	b.buffer[b.next] = entry
	b.next++
}
