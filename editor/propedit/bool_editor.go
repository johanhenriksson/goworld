package propedit

import (
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/checkbox"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

type BoolProps struct {
	Value    bool
	OnChange func(bool)
}

func BoolField(key string, title string, props BoolProps) node.T {
	return rect.New(key, rect.Props{
		Style: rect.Style{
			Layout:  style.Row{},
			Width:   style.Pct(100),
			Padding: style.RectY(4),
		},
		Children: []node.T{
			label.New("label", label.Props{
				Text: title,
				Style: label.Style{
					Grow: style.Grow(1),
				},
			}),
			checkbox.New("value", checkbox.Props{
				Checked:  props.Value,
				OnChange: props.OnChange,
			}),
		},
	})
}
