package material

import (
	"github.com/johanhenriksson/goworld/engine/uniform" // illegal import
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Pass string

const (
	Deferred = Pass("deferred")
	Forward  = Pass("forward")
)

type Standard Instance[*Descriptors]

type Descriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Lights   *descriptor.Storage[uniform.Light]
	Textures *descriptor.SamplerArray
}

type Def struct {
	Shader       string
	Pass         Pass
	VertexFormat any
	DepthTest    bool
	DepthWrite   bool
	DepthClamp   bool
	DepthFunc    core1_0.CompareOp
	Primitive    vertex.Primitive
	CullMode     vertex.CullMode
}

func (d *Def) Hash() uint64 {
	return Hash(d)
}

func FromDef(dev device.T, pool descriptor.Pool, rpass renderpass.T, def *Def) T[*Descriptors] {
	desc := &Descriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   1000,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  100,
		},
	}

	pointers := vertex.ParsePointers(def.VertexFormat)

	return New(
		dev,
		Args{
			Shader:     shader.New(dev, def.Shader),
			Pass:       rpass,
			Subpass:    renderpass.Name("main"),
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			DepthClamp: def.DepthClamp,
			DepthFunc:  def.DepthFunc,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		},
		desc)
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
