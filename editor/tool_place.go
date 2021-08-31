package editor

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type PlaceTool struct {
	object.T
	box *box.T
}

func NewPlaceTool() *PlaceTool {
	pt := &PlaceTool{
		T: object.New("PlaceTool"),
	}
	pt.box = box.Attach(pt.T, box.Args{Size: vec3.One, Color: render.Blue})
	pt.SetActive(false)
	return pt
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

func (pt *PlaceTool) Hover(editor *Editor, position, normal vec3.T) {
	pt.Transform().SetPosition(position.Add(normal.Scaled(0.5)).Floor())
}
