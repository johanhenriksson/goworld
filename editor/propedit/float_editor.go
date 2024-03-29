package propedit

import (
	"strconv"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
)

func init() {
	Register[float32](func(key, name string, prop object.GenericProp) node.T {
		return FloatField(key, name, FloatProps{
			Value:    prop.GetAny().(float32),
			OnChange: func(f float32) { prop.SetAny(f) },
		})
	})
}

type FloatProps struct {
	Label    string
	Value    float32
	OnChange func(float32)
	Validate func(float32) bool
}

func FloatField(key string, title string, props FloatProps) node.T {
	return Field(key, title, []node.T{
		Float(key, props),
	})
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

func fmtFloat(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}
