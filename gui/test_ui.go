package gui

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/label"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func TestUI(gut float32) rect.T {
	count, setCount := hooks.UseInt(0)
	onclick := func(e mouse.Event) {
		setCount(count + 1)
	}

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
			Text:    fmt.Sprintf("Clicks: %d", count),
			Size:    16.0,
			Color:   color.White,
			OnClick: onclick,
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
