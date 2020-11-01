package editor

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type EraseTool struct {
	*object.T
	box *geometry.Box
}

func NewEraseTool() *EraseTool {
	et := &EraseTool{
		T:   object.New("EraseTool"),
		box: geometry.NewBox(vec3.One, render.Red),
	}
	et.Attach(et.box)
	return et
}

func (pt *EraseTool) String() string {
	return "EraseTool"
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
	pt.Position = position.Sub(normal.Scaled(0.5)).Floor()
}
