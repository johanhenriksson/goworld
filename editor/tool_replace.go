package editor

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type ReplaceTool struct {
	*object.T
	box *geometry.Box
}

func NewReplaceTool() *ReplaceTool {
	rt := &ReplaceTool{
		T:   object.New("ReplaceTool"),
		box: geometry.NewBox(vec3.One, render.Yellow),
	}
	rt.Attach(rt.box)
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
	pt.SetPosition(position.Sub(normal.Scaled(0.5)).Floor())
}
