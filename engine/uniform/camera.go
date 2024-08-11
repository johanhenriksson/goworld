package uniform

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec4"
)

type Camera struct {
	Proj        mat4.T
	View        mat4.T
	ViewProj    mat4.T
	ProjInv     mat4.T
	ViewInv     mat4.T
	ViewProjInv mat4.T
	Eye         vec4.T
	Forward     vec4.T
	Viewport    vec2.T
	Delta       float32
	Time        float32
}

func CameraFromArgs(args draw.Args) Camera {
	return Camera{
		Proj:        args.Camera.Proj,
		View:        args.Camera.View,
		ViewProj:    args.Camera.ViewProj,
		ProjInv:     args.Camera.Proj.Invert(),
		ViewInv:     args.Camera.ViewInv,
		ViewProjInv: args.Camera.ViewProjInv,
		Eye:         vec4.Extend(args.Camera.Position, 0),
		Forward:     vec4.Extend(args.Camera.Forward, 0),
		Viewport:    vec2.NewI(args.Camera.Viewport.Width, args.Camera.Viewport.Height),

		// todo: timing values should not be part of the camera

		Delta: args.Delta,
		Time:  args.Time,
	}
}
