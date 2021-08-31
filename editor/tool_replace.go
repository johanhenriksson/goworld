package editor

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type ReplaceTool struct {
	object.T
	box *box.T
}

func NewReplaceTool() *ReplaceTool {
	rt := &ReplaceTool{
		T: object.New("ReplaceTool"),
	}
	rt.box = box.Attach(rt.T, box.Args{Size: vec3.One, Color: render.Yellow})
	rt.SetActive(false)
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
	pt.Transform().SetPosition(position.Sub(normal.Scaled(0.5)).Floor())
}
