package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

type MaterialCache cache.T[*material.Def, []Material]

// Material implements render logic for a specific material.
type Material interface {
	ID() material.ID
	Destroy()

	// Begin is called prior to drawing, once per frame.
	// Its purpose is to clear object buffers and set up per-frame data such as cameras & lighting.
	Begin(uniform.Camera, []light.T)

	// BeginGroup is called just prior to recording draw calls.
	// Its called once for each group, and may be called multiple times each frame.
	// The primary use for it is to bind the material prior to drawing each group.
	Bind(command.Recorder)

	// Draw is called for each mesh in the group.
	// Its purpose is to set up per-draw data such as object transforms and textures
	// as well as issuing the draw call.
	Draw(command.Recorder, mesh.Mesh)

	// End is called after all draw groups have been processed.
	// It runs once per frame and is primarily responsible for flushing uniform buffers.
	End()
}

func AssignMeshTextures(samplers cache.SamplerCache, msh mesh.Mesh, slots []texture.Slot) [4]uint32 {
	textureIds := [4]uint32{}
	for id, slot := range slots {
		ref := msh.Texture(slot)
		if ref != nil {
			handle, exists := samplers.TryFetch(ref)
			if exists {
				textureIds[id] = uint32(handle.ID)
			}
		}
	}
	return textureIds
}

func CameraFromArgs(args render.Args) uniform.Camera {
	return uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         vec4.Extend(args.Position, 0),
		Forward:     vec4.Extend(args.Forward, 0),
		Viewport:    vec2.NewI(args.Viewport.Width, args.Viewport.Height),
	}
}
