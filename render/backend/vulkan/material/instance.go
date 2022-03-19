package material

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
)

type Instance[D descriptor.Set] interface {
	Material() T[D]
	Descriptors() D
	Bind(command.Buffer)
}

type instance[D descriptor.Set] struct {
	material T[D]
	set      D
}

func (i *instance[D]) Material() T[D] { return i.material }
func (i *instance[D]) Descriptors() D { return i.set }

func (s *instance[D]) Bind(cmd command.Buffer) {
	s.material.Bind(cmd)
	cmd.CmdBindGraphicsDescriptor(s.material.Layout(), s.set)
}
