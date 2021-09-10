package label

import "github.com/johanhenriksson/goworld/render"

type Renderer interface {
	Draw(Label, render.Args)
}

type renderer struct {
	fontname string
	size     float32
	spacing  float32

	font *render.Font
	tex  *render.Texture
}
