package editor

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Editor base struct
type Editor struct {
	*engine.Group

	Chunk   *game.Chunk
	Camera  *engine.Camera
	Palette *PaletteWindow
	Tool    Tool

	PlaceTool   *PlaceTool
	EraseTool   *EraseTool
	SampleTool  *SampleTool
	ReplaceTool *ReplaceTool

	XPlane *geometry.Plane

	bounds  *geometry.Box
	mesh    *game.ChunkMesh
	gbuffer *render.GeometryBuffer
}

// NewEditor creates a new editor application
func NewEditor(chunk *game.Chunk, camera *engine.Camera, gbuffer *render.GeometryBuffer) *Editor {
	e := &Editor{
		Group:   engine.NewGroup(vec3.Zero, vec3.Zero),
		Chunk:   chunk,
		Camera:  camera,
		Palette: NewPaletteWindow(render.DefaultPalette),

		PlaceTool:   NewPlaceTool(),
		EraseTool:   NewEraseTool(),
		SampleTool:  NewSampleTool(),
		ReplaceTool: NewReplaceTool(),

		XPlane: geometry.NewPlane(16, render.Red),

		mesh:    game.NewChunkMesh(chunk),
		bounds:  geometry.NewBox(vec3.NewI(chunk.Sx, chunk.Sy, chunk.Sz), render.DarkGrey),
		gbuffer: gbuffer,
	}

	e.XPlane.Passes.Set(render.Geometry)
	e.XPlane.Rotation.X = -90
	e.XPlane.Position = vec3.New(8, 8, 0.01)
	e.XPlane.SetMaterial(assets.GetMaterialCached("color"))

	e.Tool = e.PlaceTool

	e.Attach(e.mesh, e.bounds, e.XPlane)

	return e
}

func (e *Editor) Collect(pass engine.DrawPass, args engine.DrawArgs) {
	e.Group.Collect(pass, args)
	engine.Collect(pass, args, e.Tool)
}

func (e *Editor) Update(dt float32) {
	e.Group.Update(dt)
	engine.Update(dt, e.Tool)

	exists, position, normal := e.cursorPositionNormal()
	if !exists {
		return
	}

	if e.Tool != nil {
		e.Tool.Hover(e, position, normal)

		// use active tool
		if mouse.Pressed(mouse.Button2) {
			e.Tool.Use(e, position, normal)
		}
	}

	// deselect tool
	if keys.Pressed(keys.Escape) {
		e.Tool = nil
	}

	// place tool
	if keys.Pressed(keys.F) {
		e.Tool = e.PlaceTool
	}

	// erase tool
	if keys.Pressed(keys.C) {
		e.Tool = e.EraseTool
	}

	// replace tool
	if keys.Pressed(keys.R) {
		e.Tool = e.ReplaceTool
	}

	// eyedropper tool
	if keys.Pressed(keys.T) {
		e.Tool = e.SampleTool
	}
}

// sample world position at current mouse coords
func (e *Editor) cursorPositionNormal() (bool, vec3.T, vec3.T) {
	depth, depthExists := e.gbuffer.SampleDepth(mouse.Position)
	if !depthExists {
		return false, vec3.Zero, vec3.Zero
	}

	viewNormal, normalExists := e.gbuffer.SampleNormal(mouse.Position)
	if !normalExists {
		return false, vec3.Zero, vec3.Zero
	}

	position := e.Camera.Unproject(vec3.Extend(
		mouse.Position.Div(e.gbuffer.Depth.Size()),
		depth,
	))

	viewInv := e.Camera.View.Invert()
	normal := viewInv.TransformDir(viewNormal)

	return true, position, normal
}

// editor components:
// - arcball camera (low prio)
// - tools
//   place voxel
//     1. palette
//     2. destination box - perhaps even ghost voxel?
//     3. placement grids
//   remove voxel
// 	   1. source box
