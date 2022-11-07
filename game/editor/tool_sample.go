package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type SampleTool struct {
	object.T
	box *box.T
}

func NewSampleTool() *SampleTool {
	st := &SampleTool{
		T: object.New("Sample"),
	}

	box.Builder(&st.box, box.Args{
		Size:  vec3.One,
		Color: color.Purple,
	}).
		Parent(st).
		Create()

	return st
}

func (pt *SampleTool) Use(editor T, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	voxel := editor.GetVoxel(int(target.X), int(target.Y), int(target.Z))
	editor.SelectColor(color.RGB8(voxel.R, voxel.G, voxel.B))

	// select placement tool
	// e.SelectTool(e.PlaceTool())
}

func (pt *SampleTool) Hover(editor T, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.box.SetActive(true)
		pt.Transform().SetPosition(p.Floor())
	} else {
		pt.box.SetActive(false)
	}
}
