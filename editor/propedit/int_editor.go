package propedit

import (
	"strconv"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
)

func init() {
	Register[int](func(key, name string, prop object.GenericProp) node.T {
		return IntegerField(key, name, IntegerProps{
			Value:    prop.GetAny().(int),
			OnChange: func(f int) { prop.SetAny(f) },
		})
	})
}

type IntegerProps struct {
	Label    string
	Value    int
	OnChange func(int)
	Validate func(int) bool
}

func IntegerField(key string, title string, props IntegerProps) node.T {
	return Field(key, title, []node.T{
		Integer(key, props),
	})
}

func Integer(key string, props IntegerProps) node.T {
	return node.Component(key, props, func(props IntegerProps) node.T {
		validate := func(int) bool { return true }
		if props.Validate != nil {
			validate = props.Validate
		}

		return String(key, StringProps{
			Label: props.Label,
			Value: fmtInteger(props.Value),
			OnChange: func(text string) {
				f, _ := strconv.ParseInt(text, 10, 64)
				props.OnChange(int(f))
			},
			Validate: func(text string) bool {
				f, err := strconv.ParseInt(text, 10, 64)
				return err == nil && validate(int(f))
			},
		})
	})
}

func fmtInteger(v int) string {
	return strconv.FormatInt(int64(v), 10)
}
