package terrain

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui/widget/icon"
	"github.com/johanhenriksson/goworld/render/color"
)

func init() {
	editor.Register(&Mesh{}, NewEditor)
}

type Editor struct {
	object.Object
	*editor.Context

	mesh *Mesh

	RaiseTool  *BrushTool
	LowerTool  *BrushTool
	SmoothTool *BrushTool
	PaintTool  *BrushTool
	LevelTool  *BrushTool
}

var _ editor.T = &Editor{}

func NewEditor(ctx *editor.Context, mesh *Mesh) *Editor {
	terrain := mesh.Tile.Map
	return object.New("TerrainEditor", &Editor{
		Context: ctx,
		mesh:    mesh,

		RaiseTool: object.Builder(NewBrushTool(terrain, NewRaiseBrush(), color.Green)).
			Active(false).
			Create(),

		LowerTool: object.Builder(NewBrushTool(terrain, NewLowerBrush(), color.Red)).
			Active(false).
			Create(),

		SmoothTool: object.Builder(NewBrushTool(terrain, &SmoothBrush{}, color.Yellow)).
			Active(false).
			Create(),

		LevelTool: object.Builder(NewBrushTool(terrain, &LevelBrush{}, color.Blue)).
			Active(false).
			Create(),

		PaintTool: object.Builder(NewBrushTool(terrain, &PaintBrush{}, color.Purple)).
			Active(false).
			Create(),
	})
}

func (e *Editor) Name() string {
	return "Terrain"
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

func (e *Editor) Actions() []editor.Action {
	return []editor.Action{
		{
			Name: "Raise",
			Icon: icon.IconArrowUpward,
			Key:  keys.F,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.RaiseTool)
			},
		},
		{
			Name: "Lower",
			Icon: icon.IconArrowDownward,
			Key:  keys.C,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.LowerTool)
			},
		},
		{
			Name: "Smooth",
			Icon: icon.IconWaves,
			Key:  keys.R,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.SmoothTool)
			},
		},
		{
			Name: "Level",
			Icon: icon.IconSyncAlt,
			Key:  keys.G,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.LevelTool)
			},
		},
		{
			Name: "Paint",
			Icon: icon.IconBrush,
			Key:  keys.T,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.PaintTool)
			},
		},
	}
}
