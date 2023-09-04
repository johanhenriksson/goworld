package terrain

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
)

func init() {
	editor.Register(&Mesh{}, NewEditor)
}

type Editor struct {
	object.Object
	*editor.Context

	mesh *Mesh
	Tile *Tile

	RaiseTool  *BrushTool
	SmoothTool *BrushTool
}

var _ editor.T = &Editor{}

func NewEditor(ctx *editor.Context, mesh *Mesh) *Editor {
	return object.New("TerrainEditor", &Editor{
		Context: ctx,
		mesh:    mesh,
		Tile:    mesh.Tile,

		RaiseTool: object.Builder(NewBrushTool(&RaiseBrush{})).
			Active(false).
			Create(),

		SmoothTool: object.Builder(NewBrushTool(&SmoothBrush{})).
			Active(false).
			Create(),
	})
}

func (e *Editor) Name() string {
	return "Chunk"
}

func (e *Editor) Target() object.Component { return e.mesh }

func (e *Editor) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.mesh.Update(scene, dt)
}

func (e *Editor) Select(ev mouse.Event) {
}

func (e *Editor) Deselect(ev mouse.Event) bool {
	return true
}

func (e *Editor) Recalculate() {
	e.mesh.Refresh()
}

func (e *Editor) Actions() []editor.Action {
	return []editor.Action{
		{
			Name: "Raise",
			Icon: icon.IconEdit,
			Key:  keys.F,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.RaiseTool)
			},
		},
		{
			Name: "Smooth",
			Icon: icon.IconEdit,
			Key:  keys.R,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.SmoothTool)
			},
		},
	}
}
