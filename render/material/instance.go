package material

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
)

type Instance[D descriptor.Set] struct {
	material *Material[D]
	set      D
}

func (i *Instance[D]) Material() *Material[D] { return i.material }
func (i *Instance[D]) Descriptors() D         { return i.set }

func (s *Instance[D]) Bind(cmd *command.Buffer) {
	// might want to move this to the command buffer instead to avoid the import
	s.material.Bind(cmd)
	cmd.CmdBindGraphicsDescriptor(s.set)
}
