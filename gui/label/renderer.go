package label

import "github.com/johanhenriksson/goworld/render"

type Renderer interface {
	Draw(render.Args, T, *Props)
}

type renderer struct {
	fontname string
	size     float32
	font     *render.Font
	tex      *render.Texture
}

func (r *renderer) Draw(args render.Args, label T, props *Props) {

}
