package gui

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/label"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type Manager interface {
	object.T
	DrawPass()
}

type manager struct {
	object.T
	scene rect.T
}

func MyCustomUI(gut float32) widget.T {
	return rect.New(
		"frame",
		&rect.Props{
			Color:  color.Hex("#000000"),
			Border: 5,
			Width:  dimension.Fixed(250),
			Height: dimension.Fixed(150),
			Layout: layout.Column{
				Padding: 5,
				Gutter:  5,
			},
		},
		label.New("title", &label.Props{
			Text:  "Hello GUI",
			Size:  16.0,
			Color: color.White,
		}),
		label.New("title2", &label.Props{
			Text:  "Hello GUI",
			Size:  16.0,
			Color: color.White,
		}),
		rect.New(
			"r1",
			&rect.Props{
				Height: dimension.Percent(150),
				Layout: layout.Row{
					Gutter:   5,
					Relative: true,
				},
			},
			rect.New("1st", &rect.Props{Color: color.Blue, Width: dimension.Fixed(1)}),
			rect.New("2nd", &rect.Props{Color: color.Green, Width: dimension.Fixed(1)}),
			rect.New("3nd", &rect.Props{Color: color.Red, Width: dimension.Fixed(2)}),
		),
		rect.New(
			"r2",
			&rect.Props{
				Height: dimension.Percent(50),
				Layout: layout.Row{
					Gutter: gut,
				},
			},
			rect.New("1st", &rect.Props{Color: color.Red}),
			rect.New("2nd", &rect.Props{Color: color.Green}),
			rect.New("3nd", &rect.Props{Color: color.Blue}),
		),
	)
}

func New() Manager {
	f := MyCustomUI(5)
	f.Move(vec2.New(400, 200))

	b := MyCustomUI(6)
	reconcile(f, b, 0)

	scene := rect.New("GUI", &rect.Props{
		Layout: layout.Absolute{},
	}, f)
	scene.Resize(vec2.New(1600, 900))

	return &manager{
		T:     object.New("GUI Manager"),
		scene: scene,
	}
}

func (m *manager) DrawPass() {
	width, height := float32(1600), float32(900)

	viewport := mat4.Orthographic(0, width, height, 0, 1000, -1000)
	scale := mat4.Ident() // todo: ui scaling
	vp := scale.Mul(&viewport)

	gl.DepthFunc(gl.LEQUAL)

	args := render.Args{
		Projection: viewport,
		View:       scale,
		VP:         viewport,
		MVP:        vp,
		Transform:  mat4.Ident(),
	}

	m.scene.Draw(args)
}

func (m *manager) MouseEvent(e mouse.Event) {
	for _, frame := range m.scene.Children() {
		if handler, ok := frame.(mouse.Handler); ok {
			handler.MouseEvent(e)
			if e.Handled() {
				break
			}

			// mouse enter
			// mouse exit
			// mouse move
			// click

			// focus
			// blur

			// keydown
			// keyup
			// keychar
		}
	}
}
