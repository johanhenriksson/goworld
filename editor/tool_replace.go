package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type ReplaceTool struct {
	object.T
	box *box.T
}

func NewReplaceTool() *ReplaceTool {
	rt := &ReplaceTool{
		T: object.New("ReplaceTool"),
	}

	box.Builder(&rt.box, box.Args{
		Size:  vec3.One,
		Color: color.Yellow,
	}).
		Parent(rt).
		Create()

	return rt
}

func (pt *ReplaceTool) String() string {
	return "ReplaceTool"
}

func (pt *ReplaceTool) Use(e *Editor, position, normal vec3.T) {
	target := position.Sub(normal.Scaled(0.5))
	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), game.NewVoxel(e.Palette.Selected))

	// recompute mesh
	e.Chunk.Light.Calculate()
	e.mesh.Compute()

	// write to disk
	go e.Chunk.Write("chunks")
}

func (pt *ReplaceTool) Hover(editor *Editor, position, normal vec3.T) {
	p := position.Sub(normal.Scaled(0.5))
	if editor.InBounds(p) {
		pt.Transform().SetPosition(p.Floor())
	}
}
