package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
)

type copyGbuffer struct {
	gbuffer GeometryBuffer
}

func NewGBufferCopyPass(gbuffer GeometryBuffer) Pass {
	return &copyGbuffer{
		gbuffer: gbuffer,
	}
}

func (p *copyGbuffer) Name() string { return "GBufferCopy" }
func (p *copyGbuffer) Destroy()     {}

func (p *copyGbuffer) Record(cmds command.Recorder, args render.Args, scene object.T) {
	p.gbuffer.RecordBufferCopy(cmds)
}
