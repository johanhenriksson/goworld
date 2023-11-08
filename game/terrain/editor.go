package terrain

import (
	"fmt"
	"log"
	"os"

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

	mesh  *Mesh
	world *World

	RaiseTool  *BrushTool
	LowerTool  *BrushTool
	SmoothTool *BrushTool
	LevelTool  *BrushTool
	NoiseTool  *BrushTool
	PaintTool  *BrushTool
}

var _ editor.T = &Editor{}

func NewEditor(ctx *editor.Context, mesh *Mesh) *Editor {
	world := object.GetInParents[*World](mesh)
	if world == nil {
		panic("mesh is not attached to a world")
	}

	mesh.Tile.Changed.Subscribe(func(t *Tile) {
		log.Println("saving tile", mesh.Parent())
		pos := mesh.Tile.Position
		path := fmt.Sprintf("assets/maps/default/tile_%d_%d", pos.X, pos.Y)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Println("failed to open tile file: %w", err)
		}
		if err := object.Save(f, mesh.Parent()); err != nil {
			log.Println("failed to save tile:", err)
		}
	})

	return object.New("TerrainEditor", &Editor{
		Context: ctx,
		mesh:    mesh,
		world:   world,

		RaiseTool: object.Builder(NewBrushTool(world.Terrain, NewRaiseBrush(), color.Green)).
			Active(false).
			Create(),

		LowerTool: object.Builder(NewBrushTool(world.Terrain, NewLowerBrush(), color.Red)).
			Active(false).
			Create(),

		SmoothTool: object.Builder(NewBrushTool(world.Terrain, &SmoothBrush{}, color.Yellow)).
			Active(false).
			Create(),

		LevelTool: object.Builder(NewBrushTool(world.Terrain, &LevelBrush{}, color.Blue)).
			Active(false).
			Create(),

		NoiseTool: object.Builder(NewBrushTool(world.Terrain, NewNoiseBrush(), color.Cyan)).
			Active(false).
			Create(),

		PaintTool: object.Builder(NewBrushTool(world.Terrain, &PaintBrush{}, color.Purple)).
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
			Icon: icon.IconSyncAlt,
			Key:  keys.R,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.SmoothTool)
			},
		},
		{
			Name: "Level",
			Icon: icon.IconTrendingFlat,
			Key:  keys.G,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.LevelTool)
			},
		},
		{
			Name: "Noise",
			Icon: icon.IconWaves,
			Key:  keys.N,
			Callback: func(mgr *editor.ToolManager) {
				mgr.UseTool(e.NoiseTool)
			},
		},

		{
			Name: "Texture 1",
			Key:  keys.Key1,
			Icon: icon.IconBrush,
			Callback: func(mgr *editor.ToolManager) {
				if mgr.Tool() != e.PaintTool {
					mgr.UseTool(e.PaintTool)
				}
				e.PaintTool.Brush.(*PaintBrush).Texture = 0
			},
		},
		{
			Name: "Texture 2",
			Key:  keys.Key2,
			Icon: icon.IconBrush,
			Callback: func(mgr *editor.ToolManager) {
				if mgr.Tool() != e.PaintTool {
					mgr.UseTool(e.PaintTool)
				}
				e.PaintTool.Brush.(*PaintBrush).Texture = 1
			},
		},
		{
			Name: "Texture 3",
			Key:  keys.Key3,
			Icon: icon.IconBrush,
			Callback: func(mgr *editor.ToolManager) {
				if mgr.Tool() != e.PaintTool {
					mgr.UseTool(e.PaintTool)
				}
				e.PaintTool.Brush.(*PaintBrush).Texture = 2
			},
		},
		{
			Name: "Texture 4",
			Key:  keys.Key4,
			Icon: icon.IconBrush,
			Callback: func(mgr *editor.ToolManager) {
				if mgr.Tool() != e.PaintTool {
					mgr.UseTool(e.PaintTool)
				}
				e.PaintTool.Brush.(*PaintBrush).Texture = 3
			},
		},
	}
}
