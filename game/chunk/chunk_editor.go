package chunk

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/window/modal"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Editor interface {
	object.Component
	editor.T

	SelectedColor() color.T
	GetVoxel(x, y, z int) voxel.T
	SetVoxel(x, y, z int, v voxel.T)
	SelectColor(color.T)
	Recalculate()
	InBounds(p vec3.T) bool
}

// Editor base struct
type edit struct {
	object.Object
	*editor.Context

	// editor target
	mesh  *Mesh
	Chunk *T
	color color.T

	PlaceTool   *PlaceTool
	EraseTool   *EraseTool
	SampleTool  *SampleTool
	ReplaceTool *ReplaceTool

	XPlane *plane.Plane
	YPlane *plane.Plane
	ZPlane *plane.Plane

	GUI  gui.Fragment
	Menu gui.Fragment

	xp, yp, zp int

	BoundingBox *box.Mesh
}

var _ editor.T = &edit{}

func init() {
	editor.Register(&Mesh{}, NewEditor)
}

// NewEditor creates a new chunk editor
func NewEditor(ctx *editor.Context, mesh *Mesh) Editor {
	chk := mesh.Chunk
	dimensions := vec3.NewI(chk.Sx, chk.Sy, chk.Sz)
	center := dimensions.Scaled(0.5)
	constructionPlaneAlpha := float32(0.33)

	e := object.New("ChunkEditor", &edit{
		Context: ctx,
		mesh:    mesh,
		Chunk:   chk,

		PlaceTool: object.Builder(NewPlaceTool()).
			Active(false).
			Create(),

		EraseTool: object.Builder(NewEraseTool()).
			Active(false).
			Create(),

		SampleTool: object.Builder(NewSampleTool()).
			Active(false).
			Create(),

		ReplaceTool: object.Builder(NewReplaceTool()).
			Active(false).
			Create(),

		color: color.Red,

		BoundingBox: box.New(box.Args{
			Size:  dimensions,
			Color: color.White,
		}),

		// X Construction Plane
		XPlane: object.Builder(plane.NewObject(plane.Args{
			Size: float32(chk.Sx),
		})).
			Position(center.WithX(0)).
			Rotation(quat.Euler(-90, 0, 90)).
			Texture(texture.Diffuse, color.Red.WithAlpha(constructionPlaneAlpha)).
			Active(false).
			Create(),

		// Y Construction Plane
		YPlane: object.Builder(plane.NewObject(plane.Args{
			Size: float32(chk.Sy),
		})).
			Position(center.WithY(0)).
			Texture(texture.Diffuse, color.Green.WithAlpha(constructionPlaneAlpha)).
			Active(false).
			Create(),

		// Z Construction Plane
		ZPlane: object.Builder(plane.NewObject(plane.Args{
			Size: float32(chk.Sz),
		})).
			Position(center.WithZ(0)).
			Rotation(quat.Euler(-90, 0, 0)).
			Texture(texture.Diffuse, color.Blue.WithAlpha(constructionPlaneAlpha)).
			Active(false).
			Create(),
	})

	e.GUI = NewGUI(e, mesh)
	object.Attach(e, e.GUI)

	e.Menu = NewMenu(e)
	object.Attach(e, e.Menu)

	return e
}

func (e *edit) Name() string {
	return "Chunk"
}

func (e *edit) Target() object.Component { return e.mesh }

func (e *edit) Select(ev mouse.Event) {
	object.Enable(e.GUI)
	object.Enable(e.Menu)
}

func (e *edit) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	object.Disable(e.Menu)
	return true
}

func (e *edit) Bounds() physics.Shape {
	return nil
}

func (e *edit) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.mesh.Update(scene, dt)
}

func (e *edit) GetVoxel(x, y, z int) voxel.T {
	return e.Chunk.At(x, y, z)
}

func (e *edit) SetVoxel(x, y, z int, v voxel.T) {
	e.Chunk.Set(x, y, z, v)
}

func (e *edit) SelectColor(c color.T) {
	e.color = c
}

func (e *edit) SelectedColor() color.T {
	return e.color
}

func (e *edit) InBounds(p vec3.T) bool {
	p = p.Floor()
	outside := p.X < 0 || p.Y < 0 || p.Z < 0 || int(p.X) >= e.Chunk.Sx || int(p.Y) >= e.Chunk.Sy || int(p.Z) >= e.Chunk.Sz
	return !outside
}

func (e *edit) clearChunk() {
	e.Chunk.Clear()
	e.Recalculate()
}

func (e *edit) KeyEvent(ev keys.Event) {
	// toggle construction planes
	if keys.PressedMods(ev, keys.X, keys.Ctrl) {
		object.Toggle(e.XPlane, !e.XPlane.Enabled())
	}
	if keys.PressedMods(ev, keys.Y, keys.Ctrl) {
		object.Toggle(e.YPlane, !e.YPlane.Enabled())
	}
	if keys.PressedMods(ev, keys.Z, keys.Ctrl) {
		object.Toggle(e.ZPlane, !e.ZPlane.Enabled())
	}

	m := 1
	if ev.Modifier(keys.Shift) {
		m = -1
	}

	if keys.Pressed(ev, keys.X) && e.XPlane.Active() {
		e.xp = (e.xp + e.Chunk.Sx + m + 1) % (e.Chunk.Sx + 1)
		p := e.XPlane.Transform().Position().WithX(float32(e.xp))
		e.XPlane.Transform().SetPosition(p)
	}

	if keys.Pressed(ev, keys.Y) && e.YPlane.Active() {
		e.yp = (e.yp + e.Chunk.Sy + m + 1) % (e.Chunk.Sy + 1)
		p := e.YPlane.Transform().Position().WithY(float32(e.yp))
		e.YPlane.Transform().SetPosition(p)
	}

	if keys.Pressed(ev, keys.Z) && e.ZPlane.Active() {
		e.zp = (e.zp + e.Chunk.Sz + m + 1) % (e.Chunk.Sz + 1)
		p := e.ZPlane.Transform().Position().WithZ(float32(e.zp))
		e.ZPlane.Transform().SetPosition(p)
	}
}

func (e *edit) Recalculate() {
	e.Chunk.Light.Calculate()
	e.mesh.Refresh()
}

func (e *edit) Actions() []editor.Action {
	return []editor.Action{
		{
			Name:     "Place",
			Key:      keys.F,
			Callback: func(mgr editor.ToolManager) { mgr.SelectTool(e.PlaceTool) },
		},
		{
			Name:     "Erase",
			Key:      keys.C,
			Callback: func(mgr editor.ToolManager) { mgr.SelectTool(e.EraseTool) },
		},
		{
			Name:     "Replace",
			Key:      keys.R,
			Callback: func(mgr editor.ToolManager) { mgr.SelectTool(e.ReplaceTool) },
		},
		{
			Name: "Sample",
			Key:  keys.T,
			Callback: func(mgr editor.ToolManager) {
				previousTool := mgr.Tool()
				e.SampleTool.Reselect = func() { mgr.SelectTool(previousTool) }
				mgr.SelectTool(e.SampleTool)
			},
		},
		{
			Name:     "Clear",
			Key:      keys.N,
			Modifier: keys.Ctrl,
			Callback: func(mgr editor.ToolManager) { e.clearChunk() },
		},
		{
			Name:     "Save",
			Key:      keys.S,
			Modifier: keys.Ctrl,
			Callback: func(mgr editor.ToolManager) { e.saveChunkDialog() },
		},
	}
}

func (e *edit) saveChunkDialog() {
	var saveDialog gui.Fragment
	saveDialog = gui.NewFragment(gui.FragmentArgs{
		Slot:     "gui",
		Position: gui.FragmentFirst,
		Render: func() node.T {
			return modal.NewInput("modal test", modal.InputProps{
				Title:   "Save as...",
				Message: "Enter filename:",
				OnClose: func() {
					object.Detach(saveDialog)
				},
				OnAccept: func(input string) {
					log.Println("input:", input)
				},
			})
		},
	})
	object.Attach(e, saveDialog)
}
