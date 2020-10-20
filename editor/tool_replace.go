package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type ReplaceTool struct {
	box *geometry.Box
}

func NewReplaceTool() *ReplaceTool {
	return &ReplaceTool{
		box: geometry.NewBox(vec3.One, render.Yellow),
	}
}

func (pt *ReplaceTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.NewVoxel(e.Palette.Selected))

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *ReplaceTool) Update(editor *Editor, dt float32, position, normal vec3.T) {
	pt.box.Position = position.Sub(normal.Scaled(0.5)).Floor()
	engine.Update(dt, pt.box)
}

func (pt *ReplaceTool) Draw(editor *Editor, args engine.DrawArgs) {
	engine.Draw(args, pt.box)
}
