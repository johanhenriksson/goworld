package editor

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/voxel"
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
		T: object.New("Place"),
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
	x, y, z := int(target.X), int(target.Y), int(target.Z)

	if editor.GetVoxel(x, y, z) != voxel.Empty {
		return
	}

	clr := editor.SelectedColor()
	log.Println("place", clr, "at", x, y, z)
	editor.SetVoxel(x, y, z, voxel.New(clr))

	// recompute mesh
	editor.Recalculate()
}

func (pt *PlaceTool) Hover(editor T, position, normal vec3.T) {
	p := position.Add(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}
