package editor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type EraseTool struct {
	box *geometry.ColorCube
}

func NewEraseTool() *EraseTool {
	return &EraseTool{
		box: geometry.NewColorCube(render.Red, 1.05),
	}
}

func (pt *EraseTool) Use(e *Editor, position, normal vec3.T) {
	fmt.Println("(tool) Erase at", position)
	target := position.Sub(normal.Scaled(0.5))
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.EmptyVoxel)

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *EraseTool) Update(editor *Editor, dt float32, position, normal vec3.T) {
	pt.box.Position = position.Sub(normal.Scaled(0.5)).Floor().Add(vec3.One.Scaled(0.5))
	pt.box.Update(dt)
}

func (pt *EraseTool) Draw(editor *Editor, args engine.DrawArgs) {
	pt.box.Draw(args)
}
