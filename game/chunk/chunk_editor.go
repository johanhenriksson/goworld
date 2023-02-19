package chunk

import (
	"log"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/geometry/box"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/window/modal"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type Editor interface {
	editor.T

	SelectedColor() color.T
	GetVoxel(x, y, z int) voxel.T
	SetVoxel(x, y, z int, v voxel.T)
	SelectColor(color.T)
	Recalculate()
	InBounds(p vec3.T) bool

	CursorPositionNormal(cursor vec2.T) (bool, vec3.T, vec3.T)
}

// Editor base struct
type edit struct {
	object.T

	// editor target
	mesh *Mesh

	Chunk  *T
	Camera camera.T

	color color.T

	PlaceTool   *PlaceTool
	EraseTool   *EraseTool
	SampleTool  *SampleTool
	ReplaceTool *ReplaceTool

	XPlane *plane.T
	YPlane *plane.T
	ZPlane *plane.T

	xp, yp, zp int

	Bounds *box.T
	render renderer.T
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

	e := object.New(&edit{
		mesh:   mesh,
		Chunk:  chk,
		Camera: ctx.Camera,

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

		render: ctx.Render,
		color:  color.Red,

		Bounds: box.New(box.Args{
			Size:  dimensions,
			Color: color.White,
		}),

		// X Construction Plane
		XPlane: object.Builder(plane.New(plane.Args{
			Size:  float32(chk.Sx),
			Color: color.Red.WithAlpha(0.25),
		})).
			Position(center.WithX(0)).
			Rotation(vec3.New(-90, 0, 90)).
			Active(false).
			Create(),

		// Y Construction Plane
		YPlane: object.Builder(plane.New(plane.Args{
			Size:  float32(chk.Sy),
			Color: color.Green.WithAlpha(0.25),
		})).
			Position(center.WithY(0)).
			Active(false).
			Create(),

		// Z Construction Plane
		ZPlane: object.Builder(plane.New(plane.Args{
			Size:  float32(chk.Sz),
			Color: color.Blue.WithAlpha(0.25),
		})).
			Position(center.WithZ(0)).
			Rotation(vec3.New(-90, 0, 0)).
			Active(false).
			Create(),
	})

	object.Attach(e, NewGUI(e))
	object.Attach(e, NewMenu(e))

	return e
}

func (e *edit) Name() string {
	return "Editor"
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

// sample world position at current mouse coords
func (e *edit) CursorPositionNormal(cursor vec2.T) (bool, vec3.T, vec3.T) {
	viewPosition, positionExists := e.render.GBuffer().SamplePosition(cursor)
	if !positionExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewNormal, normalExists := e.render.GBuffer().SampleNormal(cursor)
	if !normalExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewInv := e.Camera.ViewInv()
	normal := viewInv.TransformDir(viewNormal)
	position := viewInv.TransformPoint(viewPosition)

	// transform world coords into object space
	position = e.Transform().Unproject(position)

	return true, position, normal
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
		e.XPlane.SetActive(!e.XPlane.Active())
	}
	if keys.PressedMods(ev, keys.Y, keys.Ctrl) {
		e.YPlane.SetActive(!e.YPlane.Active())
	}
	if keys.PressedMods(ev, keys.Z, keys.Ctrl) {
		e.ZPlane.SetActive(!e.ZPlane.Active())
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
	e.mesh.Compute()
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
