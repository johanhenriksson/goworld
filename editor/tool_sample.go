package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type SampleTool struct {
	object.T
	box *box.T
}

func NewSampleTool() *SampleTool {
	st := &SampleTool{
		T: object.New("SampleTool"),
	}

	box.Builder(&st.box, box.Args{
		Size:  vec3.One,
		Color: render.Purple,
	}).
		Parent(st).
		Create()

	return st
}

func (pt *SampleTool) String() string {
	return "SampleTool"
}

func (pt *SampleTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	voxel := e.Chunk.At(int(target.X), int(target.Y), int(target.Z))
	e.Palette.Selected = render.Color4(float32(voxel.R)/255, float32(voxel.G)/255, float32(voxel.B)/255, 1)

	// select placement tool
	e.SelectTool(e.PlaceTool)
}

func (pt *SampleTool) Hover(editor *Editor, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.box.SetActive(true)
		pt.Transform().SetPosition(p.Floor())
	} else {
		pt.box.SetActive(false)
	}
}
