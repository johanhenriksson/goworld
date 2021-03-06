package editor

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type EraseTool struct {
	*object.T
	box *box.T
}

func NewEraseTool() *EraseTool {
	et := &EraseTool{
		T: object.New("EraseTool"),
	}
	et.box = box.Attach(et.T, box.Args{Size: vec3.One, Color: render.Red})
	et.SetActive(false)
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
	// parent actually refers to the editor right now
	// tools should be attached to their own object
	// they could potentially share positioning logic
	pt.SetPosition(position.Sub(normal.Scaled(0.5)).Floor())
}
