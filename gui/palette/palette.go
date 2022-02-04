package palette

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/label"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	rect.T
}

type Props struct {
	Palette color.Palette
	OnPick  func(color.T)
}

func Map[T any, S any](items []T, transform func(int, T) S) []S {
	output := make([]S, len(items))
	for i, item := range items {
		output[i] = transform(i, item)
	}
	return output
}

func Chunks[T any](slice []T, size int) [][]T {
	count := len(slice) / size
	chunks := make([][]T, 0, count)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func New(key string, props *Props) node.T {
	return node.Component(key, props, nil, func(props *Props) node.T {
		perRow := 5

		selected, setSelected := hooks.UseState(props.Palette[0])

		colors := Map(props.Palette, func(i int, c color.T) node.T {
			return rect.New(fmt.Sprintf("color%d", i), &rect.Props{
				Color: c,
				OnClick: func(e mouse.Event) {
					setSelected(c)
					if props.OnPick != nil {
						props.OnPick(c)
					}
				},
			})
		})

		rows := Map(Chunks(colors, perRow), func(i int, colors []node.T) node.T {
			return rect.New(fmt.Sprintf("row%d", i), &rect.Props{
				Layout: layout.Row{
					Padding: 1,
					Gutter:  2,
				},
				Children: colors,
			})
		})

		return rect.New("window", &rect.Props{
			Border: 3.0,
			Color:  color.Black.WithAlpha(0.8),
			Width:  dimension.Fixed(140),
			Height: dimension.Fixed(230),
			Layout: layout.Column{
				Padding: 4,
			},
			Children: []node.T{
				label.New("title", &label.Props{
					Text:  "Palette",
					Color: color.White,
					Size:  16,
				}),
				rect.New("selected", &rect.Props{
					Layout: layout.Row{},
					Height: dimension.Fixed(16),
					Children: []node.T{
						label.New("selected", &label.Props{
							Text:  "Selected",
							Color: color.White,
						}),
						rect.New("preview", &rect.Props{
							Color:  selected,
							Width:  dimension.Fixed(20),
							Height: dimension.Fixed(10),
						}),
					},
				}),
				rect.New("grid", &rect.Props{
					Height:   dimension.Fixed(200),
					Children: rows,
				}),
			},
		})
	})
}
