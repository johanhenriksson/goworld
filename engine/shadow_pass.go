package engine

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type ShadowPass struct {
	Output    *render.Texture
	Material  *render.Material
	Width     int
	Height    int
	shadowmap *render.FrameBuffer
}

func NewShadowPass(input *render.GeometryBuffer) *ShadowPass {
	mat := assets.GetMaterial("color_geometry")

	size := 4096
	fbo := render.CreateFrameBuffer(int32(size), int32(size))
	fbo.ClearColor = render.Color4(1, 1, 1, 1)
	texture := fbo.AttachBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT)

	p := &ShadowPass{
		Material:  mat,
		shadowmap: fbo,
		Output:    texture,
		Width:     size,
		Height:    size,
	}
	return p
}

func (sp *ShadowPass) DrawPass(scene *Scene, light *Light) {
	if light.Type != DirectionalLight {
		// only directional lights support shadows atm
		return
	}

	/* bind shadow map depth render target */
	sp.shadowmap.Bind()
	sp.shadowmap.Clear()

	/* use shadow pass shader - usually the standard geometry shader */
	sp.Material.Use()

	gl.DepthMask(true)

	/* compute world to lightspace (light view projection) matrix */
	// todo: move to light object
	p := light.Projection
	v := mgl.LookAtV(light.Position, mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})
	vp := p.Mul4(v)

	args := render.DrawArgs{
		Projection: p,
		View:       v,
		VP:         vp,
		MVP:        vp,
		Transform:  mgl.Ident4(),

		Pass:   render.GeometryPass,
		Shader: sp.Material.Shader,
	}
	scene.DrawCall(args)
	//scene.Draw("geometry", sp.Material.Shader)

	sp.shadowmap.Unbind()
}
