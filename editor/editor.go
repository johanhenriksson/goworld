package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Editor base struct
type Editor struct {
	Chunk   *game.Chunk
	Camera  *engine.Camera
	Palette *PaletteWindow
	Tool    Tool

	PlaceTool  *PlaceTool
	EraseTool  *EraseTool
	SampleTool *SampleTool

	mesh    *game.ChunkMesh
	gbuffer *render.GeometryBuffer
}

// NewEditor creates a new editor application
func NewEditor(chunk *game.Chunk, camera *engine.Camera, gbuffer *render.GeometryBuffer) *Editor {
	editor := &Editor{
		Chunk:   chunk,
		Camera:  camera,
		Palette: NewPaletteWindow(render.DefaultPalette),

		PlaceTool:  NewPlaceTool(),
		EraseTool:  NewEraseTool(),
		SampleTool: NewSampleTool(),

		mesh:    game.NewChunkMesh(chunk),
		gbuffer: gbuffer,
	}
	editor.Tool = editor.PlaceTool
	return editor
}

func (e *Editor) Draw(args engine.DrawArgs) {
	e.mesh.Draw(args)

	if e.Tool != nil {
		e.Tool.Draw(e, args)
	}
}

func (e *Editor) Update(dt float32) {
	e.mesh.Update(dt)

	exists, position, normal := e.cursorPositionNormal()
	if !exists {
		return
	}

	if e.Tool != nil {
		e.Tool.Update(e, dt, position, normal)

		// use active tool
		if mouse.Pressed(mouse.Button2) {
			e.Tool.Use(e, position, normal)
		}
	}

	// if keys.Released(keys.R) {
	// 	// replace voxel
	// 	fmt.Println("Replace at", position)
	// 	target := position.Sub(normal.Scaled(0.5))
	// 	e.Chunk.Set(int(target.X), int(target.Y), int(target.Z), selected)

	// 	// recompute mesh
	// 	e.Chunk.Light.Calculate()
	// 	e.mesh.Compute()

	// 	// write to disk
	// 	go e.Chunk.Write("chunks")
	// }

	// place tool
	if keys.Pressed(keys.F) {
		e.Tool = e.PlaceTool
	}

	// erase tool
	if keys.Pressed(keys.C) {
		e.Tool = e.EraseTool
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
