package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type SampleTool struct {
	box *geometry.Box
}

func NewSampleTool() *SampleTool {
	return &SampleTool{
		box: geometry.NewBox(vec3.One, render.Purple),
	}
}

func (pt *SampleTool) Name() string {
	return "SampleTool"
}

func (pt *SampleTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	voxel := e.Chunk.At(int(target.X), int(target.Y), int(target.Z))
	e.Palette.Selected = render.Color4(float32(voxel.R)/255, float32(voxel.G)/255, float32(voxel.B)/255, 1)

	// select placement tool
	e.Tool = e.PlaceTool
}

func (pt *SampleTool) Hover(editor *Editor, position, normal vec3.T) {
	pt.box.Position = position.Sub(normal.Scaled(0.5)).Floor()
}

func (pt *SampleTool) Update(dt float32) {
	engine.Update(dt, pt.box)
}

func (pt *SampleTool) Collect(pass engine.DrawPass, args engine.DrawArgs) {
	engine.Collect(pass, args, pt.box)
}
