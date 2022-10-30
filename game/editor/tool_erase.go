package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type EraseTool struct {
	object.T
	box *box.T
}

func NewEraseTool() *EraseTool {
	et := &EraseTool{
		T: object.New("EraseTool"),
	}

	box.Builder(&et.box, box.Args{
		Size:  vec3.One,
		Color: color.Red,
	}).
		Parent(et).
		Create()

	return et
}

func (pt *EraseTool) String() string {
	return "EraseTool"
}

func (pt *EraseTool) Use(editor T, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	editor.SetVoxel(int(target.X), int(target.Y), int(target.Z), voxel.Empty)

	// recompute mesh
	editor.Recalculate()
}

func (pt *EraseTool) Hover(editor T, position, normal vec3.T) {
	// parent actually refers to the editor right now
	// tools should be attached to their own object
	// they could potentially share positioning logic
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}
