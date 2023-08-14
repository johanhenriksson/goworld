package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type SpritePass struct {
	app vulkan.App
}

func NewSpritePass(app vulkan.App) *SpritePass {
	// this is basically a specialized forward pass
	return &SpritePass{
		app: app,
	}
}

func (p *SpritePass) Name() string { return "Sprites" }

func (p *SpritePass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
}

func (p *SpritePass) Destroy() {
}
