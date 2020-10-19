package editor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type PlaceTool struct {
	box *geometry.ColorCube
}

func NewPlaceTool() *PlaceTool {
	return &PlaceTool{
		box: geometry.NewColorCube(render.Blue, 1),
	}
}

func (pt *PlaceTool) Use(e *Editor, position, normal vec3.T) {
	fmt.Println("(tool) Place at", position)
	target := position.Add(normal.Scaled(0.5))
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.NewVoxel(e.Palette.Selected))

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *PlaceTool) Update(editor *Editor, dt float32, position, normal vec3.T) {
	pt.box.Position = position.Add(normal.Scaled(0.5)).Floor().Add(vec3.One.Scaled(0.5))
	pt.box.Update(dt)
}

func (pt *PlaceTool) Draw(editor *Editor, args engine.DrawArgs) {
	pt.box.Draw(args)
}
