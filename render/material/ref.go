package material

import (
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Pass string

const (
	Deferred = Pass("deferred")
	Forward  = Pass("forward")
)

type Ref interface {
	Key() string
	Version() int

	Shader() string
	Pass() Pass
	VertexFormat() any
	DepthTest() bool
	DepthWrite() bool
	DepthClamp() bool
	DepthFunc() core1_0.CompareOp
	Primitive() vertex.Primitive
	CullMode() vertex.CullMode
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
