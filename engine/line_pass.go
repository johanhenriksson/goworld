package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render"
)

// LinePass draws line geometry
type LinePass struct {
	meshes cache.Meshes
}

// NewLinePass sets up a line geometry pass.
func NewLinePass(meshes cache.Meshes) *LinePass {
	return &LinePass{
		meshes: meshes,
	}
}

// DrawPass executes the line pass
func (p *LinePass) Draw(args render.Args, scene object.T) {
	render.BindScreenBuffer()
	render.SetViewport(render.Viewport{
		Width:  args.Viewport.Width,
		Height: args.Viewport.Height,
	})

	objects := query.New[mesh.T]().Where(isDrawLines).Collect(scene)
	for _, mesh := range objects {
		p.DrawLines(args, mesh)
	}
}

func (p *LinePass) DrawLines(args render.Args, m mesh.T) error {
	args = args.Apply(m.Transform().World())
	mat := m.Material()

	if err := mat.Use(); err != nil {
		return fmt.Errorf("failed to assign material %s in mesh %s: %w", mat.Name(), m.Name(), err)
	}

	mat.Mat4("mvp", args.MVP)

	drawable := p.meshes.Fetch(m.Mesh(), mat)
	return drawable.Draw()
}

func isDrawLines(m mesh.T) bool {
	return m.Mode() == mesh.Lines
}
