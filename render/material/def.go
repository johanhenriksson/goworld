package material

import (
	"encoding/gob"
	"strconv"

	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type ID uint64

type Pass string

const (
	Deferred = Pass("deferred")
	Forward  = Pass("forward")
)

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
	Transparent  bool

	id ID
}

func init() {
	gob.Register(&Def{})
}

func (d *Def) Hash() ID {
	if d == nil {
		return 0
	}
	if d.id == 0 {
		// cache the hash
		// todo: it might be a problem that this wont ever be invalidated
		d.id = Hash(d)
	}
	return d.id
}

func (d *Def) Key() string {
	return strconv.FormatUint(uint64(d.Hash()), 16)
}

func (d *Def) Version() int {
	return 1
}

func (d *Def) LoadMaterial(fs.Filesystem) *Def {
	return d
}

func Hash(def *Def) ID {
	if def == nil {
		return 0
	}
	hash, err := hashstructure.Hash(*def, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return ID(hash)
}
