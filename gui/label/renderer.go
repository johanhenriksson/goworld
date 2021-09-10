package label

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

type Renderer interface {
	Draw(render.Args, T, *Props)
}

type renderer struct {
	text string
	font *render.Font
	tex  *render.Texture
}

func (r *renderer) Draw(args render.Args, label T, props *Props) {
	font := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", props.Size, 1.0)

	if r.text != props.Text || font != r.font {
		// (re)create opengl texture
		bounds := r.font.Measure(props.Text)
		r.tex = render.CreateTexture(int(bounds.X), int(bounds.Y))
		r.text = props.Text
		r.font = font
	}

}
