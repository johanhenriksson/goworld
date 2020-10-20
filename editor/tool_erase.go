package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type EraseTool struct {
	box *geometry.Box
}

func NewEraseTool() *EraseTool {
	return &EraseTool{
		box: geometry.NewBox(vec3.One, render.Red),
	}
}

func (pt *EraseTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.EmptyVoxel)

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *EraseTool) Hover(editor *Editor, position, normal vec3.T) {
	pt.box.Position = position.Sub(normal.Scaled(0.5)).Floor()
}

func (pt *EraseTool) Update(dt float32) {
	engine.Update(dt, pt.box)
}

func (pt *EraseTool) Collect(pass engine.DrawPass, args engine.DrawArgs) {
	pt.box.Collect(pass, args)
}
