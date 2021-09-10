package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type PlaceTool struct {
	object.T
	box *box.T
}

func NewPlaceTool() *PlaceTool {
	pt := &PlaceTool{
		T: object.New("PlaceTool"),
	}

	box.Builder(&pt.box, box.Args{
		Size:  vec3.One,
		Color: color.Blue,
	}).
		Parent(pt).
		Create()

	return pt
}

func (pt *PlaceTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Add(normal.Scaled(0.5))
	if e.Chunk.At(int(target.X), int(target.Y), int(target.Z)) != game.EmptyVoxel {
		return
	}
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.NewVoxel(e.Palette.Selected))

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *PlaceTool) Hover(editor *Editor, position, normal vec3.T) {
	p := position.Add(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}
