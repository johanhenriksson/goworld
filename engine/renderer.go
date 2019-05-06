package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	//mgl "github.com/go-gl/mathgl/mgl32"
)

type RenderPass interface {
	DrawPass(*Scene)
}

/* Maps names to Render Passes */
type PassMap map[string]RenderPass

/* Renderer - Holds references to the Scene Graph and is
 * responsible for executing render passes in order */
type Renderer struct {
	Passes   []RenderPass
	Scene    *Scene
	Width    int32
	Height   int32
	pass_map PassMap
}

/* Instantiate a new renderer. Also sets up basic OpenGL settings */
func NewRenderer(width, height int32, scene *Scene) *Renderer {
	r := &Renderer{
		Width:    width,
		Height:   height,
		Scene:    scene,
		Passes:   []RenderPass{},
		pass_map: make(PassMap),
	}
	return r
}

/* Append a new render pass */
func (r *Renderer) Append(name string, pass RenderPass) {
	r.Passes = append(r.Passes, pass)
	if len(name) > 0 {
		r.pass_map[name] = pass
	}
}

/* Get render pass by name */
func (r *Renderer) Get(name string) RenderPass {
	return r.pass_map[name]
}

/* Clears all render passes */
func (r *Renderer) Reset() {
	r.Passes = []RenderPass{}
	r.pass_map = make(PassMap)
}

func (r *Renderer) Draw() {
	/* Clear screen */
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	/* Enable blending */
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	/* Depth test */
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for _, pass := range r.Passes {
		pass.DrawPass(r.Scene)
	}
}

func (r *Renderer) Update(dt float32) {
	r.Scene.Update(dt)
}
