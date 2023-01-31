package material

import (
	"github.com/johanhenriksson/goworld/engine/renderer/uniform" // illegal import
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/mitchellh/hashstructure/v2"

	vk "github.com/vulkan-go/vulkan"
)

type Ref interface {
	cache.Key
}

type Standard Instance[*Descriptors]

type Descriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type Def struct {
	Shader       string
	Subpass      renderpass.Name
	VertexFormat any
	DepthTest    bool
	DepthWrite   bool
	Primitive    vertex.Primitive
}

func (d *Def) Hash() uint64 {
	return Hash(d)
}

func FromDef(dev device.T, pool descriptor.Pool, rpass renderpass.T, def *Def) Standard {
	desc := &Descriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: vk.ShaderStageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: vk.ShaderStageAll,
			Size:   100,
		},
		Textures: &descriptor.SamplerArray{
			Stages: vk.ShaderStageFragmentBit,
			Count:  100,
		},
	}

	pointers := vertex.ParsePointers(def.VertexFormat)

	return New(
		dev,
		Args{
			Shader:     shader.New(dev, def.Shader),
			Pass:       rpass,
			Subpass:    def.Subpass,
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			Primitive:  def.Primitive,
		},
		desc).Instantiate(pool)
}

func Hash(def *Def) uint64 {
	if def == nil {
		return 0
	}
	hash, err := hashstructure.Hash(*def, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
}
