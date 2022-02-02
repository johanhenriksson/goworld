package palette

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/label"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render/color"
)

type Props struct {
	Palette color.Palette
	OnPick  func(color.T)
}

func New(key string, props *Props) rect.T {
	perRow := 5
	rows := make([]widget.T, 0, len(props.Palette)/perRow)

	selected, setSelected := hooks.UseColor(props.Palette[0])

	items := make([]widget.T, 0, perRow)
	for i, clr := range props.Palette {
		index := i + 1
		col := index % perRow

		pickColor := clr
		items = append(items, rect.New(fmt.Sprintf("col%d", col), &rect.Props{
			Color: clr,
			OnClick: func(e mouse.Event) {
				setSelected(pickColor)
				if props.OnPick != nil {
					props.OnPick(pickColor)
				}
			},
		}))

		if index%perRow == 0 {
			row := index / perRow
			rows = append(rows, rect.New(fmt.Sprintf("row%d", row), &rect.Props{
				Layout: layout.Row{
					Padding: 1,
					Gutter:  2,
				},
			}, items...))
			items = make([]widget.T, 0, perRow)
		}
	}

	if len(items) > 0 {
		row := len(rows)
		rows = append(rows, rect.New(fmt.Sprintf("row%d", row), &rect.Props{
			Layout: layout.Row{
				Padding: 1,
				Gutter:  2,
			},
		}, items...))
	}

	grid := rect.New("grid", &rect.Props{
		Height: dimension.Fixed(200),
	}, rows...)

	return rect.New(
		key,
		&rect.Props{
			Border: 3.0,
			Color:  color.Black.WithAlpha(0.8),
			Width:  dimension.Fixed(140),
			Height: dimension.Fixed(230),
			Layout: layout.Column{
				Padding: 4,
			},
		},
		label.New("title", &label.Props{
			Text:  "Palette",
			Color: color.White,
			Size:  16,
		}),
		rect.New(
			"selected",
			&rect.Props{
				Layout: layout.Row{},
				Height: dimension.Fixed(16),
			},
			label.New("selected", &label.Props{
				Text:  "Selected",
				Color: color.White,
			}),
			rect.New(
				"preview",
				&rect.Props{
					Color:  selected,
					Width:  dimension.Fixed(20),
					Height: dimension.Fixed(10),
				},
			),
		),
		grid,
	)
}
