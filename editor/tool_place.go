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

func (pt *PlaceTool) Use(editor T, position, normal vec3.T) {
	target := position.Add(normal.Scaled(0.5))
	if editor.GetVoxel(int(target.X), int(target.Y), int(target.Z)) != game.EmptyVoxel {
		return
	}
	editor.SetVoxel(int(target.X), int(target.Y), int(target.Z), game.NewVoxel(editor.SelectedColor()))

	// recompute mesh
	editor.Recalculate()
}

func (pt *PlaceTool) Hover(editor T, position, normal vec3.T) {
	p := position.Add(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}
