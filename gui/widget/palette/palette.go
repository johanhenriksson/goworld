package palette

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

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

func New(key string, props Props) node.T {
	return node.Component(key, props, nil, render)
}

func render(props Props) node.T {
	perRow := 5

	selected, setSelected := hooks.UseState(props.Palette[0])

	colors := Map(props.Palette, func(i int, c color.T) node.T {
		return rect.New(fmt.Sprintf("color%d", i), rect.Props{
			Style: Sheet{
				Color:  c,
				Grow:   Grow(0),
				Shrink: Shrink(1),
				Basis:  Pct(16),
				Height: Px(20),
				Margin: Px(2),
			},
			OnClick: func(e mouse.Event) {
				setSelected(c)
				if props.OnPick != nil {
					props.OnPick(c)
				}
			},
		})
	})

	rows := Map(Chunks(colors, perRow), func(i int, colors []node.T) node.T {
		return rect.New(fmt.Sprintf("row%d", i), rect.Props{
			Style: Sheet{
				Width:  Pct(100),
				Layout: Row{},
			},
			Children: colors,
		})
	})

	return rect.New("window", rect.Props{
		Style: Sheet{
			Color:   color.Black.WithAlpha(0.9),
			Padding: Px(4),
			Layout:  Column{},
		},
		Children: []node.T{
			label.New("title", label.Props{
				Text: "Palette",
				Style: Sheet{
					Color: color.White,
					Font: style.Font{
						Name: "fonts/SourceCodeProRegular.ttf",
						Size: 16,
					},
				},
			}),
			rect.New("selected", rect.Props{
				Style: Sheet{
					Layout:   Row{},
					MaxWidth: Pct(100),
				},
				Children: []node.T{
					label.New("selected", label.Props{
						Text: "Selected",
						Style: Sheet{
							Color: color.White,
							Basis: Pct(80),
							Grow:  Grow(1),
						},
					}),
					rect.New("preview", rect.Props{
						Style: Sheet{
							Color:  selected,
							Grow:   Grow(1),
							Shrink: Shrink(1),
							Basis:  Px(20),
							Height: Px(20),
						},
					}),
				},
			}),
			rect.New("grid", rect.Props{
				Children: rows,
			}),
		},
	})
}
