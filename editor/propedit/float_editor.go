package propedit

import (
	"strconv"

	"github.com/johanhenriksson/goworld/gui/node"
)

type FloatProps struct {
	Label    string
	Value    float32
	OnChange func(float32)
	Validate func(float32) bool
}

func fmtFloat(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func Float(key string, props FloatProps) node.T {
	return node.Component(key, props, func(props FloatProps) node.T {
		validate := func(float32) bool { return true }
		if props.Validate != nil {
			validate = props.Validate
		}

		return String(key, StringProps{
			Label: props.Label,
			Value: fmtFloat(props.Value),
			OnChange: func(text string) {
				f, _ := strconv.ParseFloat(text, 32)
				props.OnChange(float32(f))
			},
			Validate: func(text string) bool {
				f, err := strconv.ParseFloat(text, 32)
				return err == nil && validate(float32(f))
			},
		})
	})
}
