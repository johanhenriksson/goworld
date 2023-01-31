package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type SampleTool struct {
	object.T
	Box *box.T
}

func NewSampleTool() *SampleTool {
	return object.New(&SampleTool{
		Box: box.New(box.Args{
			Size:  vec3.One,
			Color: color.Purple,
		}),
	})
}

func (pt *SampleTool) Use(editor Voxels, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	voxel := editor.GetVoxel(int(target.X), int(target.Y), int(target.Z))
	editor.SelectColor(color.RGB8(voxel.R, voxel.G, voxel.B))

	// select placement tool
	// e.SelectTool(e.PlaceTool())
}

func (pt *SampleTool) Hover(editor Voxels, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Box.SetActive(true)
		pt.Transform().SetPosition(p.Floor())
	} else {
		pt.Box.SetActive(false)
	}
}
