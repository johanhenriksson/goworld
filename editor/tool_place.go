package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type PlaceTool struct {
	box *geometry.Box
}

func NewPlaceTool() *PlaceTool {
	return &PlaceTool{
		box: geometry.NewBox(vec3.One, render.Blue),
	}
}

func (pt *PlaceTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Add(normal.Scaled(0.5))
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.NewVoxel(e.Palette.Selected))

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *PlaceTool) Update(editor *Editor, dt float32, position, normal vec3.T) {
	pt.box.Position = position.Add(normal.Scaled(0.5)).Floor()
	engine.Update(dt, pt.box)
}

func (pt *PlaceTool) Draw(editor *Editor, args engine.DrawArgs) {
	engine.Draw(args, pt.box)
}
