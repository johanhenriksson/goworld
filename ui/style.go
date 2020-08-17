package ui

import (
	"fmt"
	"github.com/johanhenriksson/goworld/render"
)

// ErrIllegalCast occurs when attempting to cast a style variable to an incompatible type
var ErrIllegalCast = fmt.Errorf("illegal cast")

// NoStyle is the empty style sheet
var NoStyle = Style{}

// Style holds UI styling information
type Style map[string]Variable

// Color returns a color value from the styles
func (s Style) Color(name string, def render.Color) render.Color {
	if value, set := s[name]; set {
		color, err := value.Color()
		if err != nil {
			panic(err)
		}
		return color
	}
	return def
}

// Float returns a float value from the styles
func (s Style) Float(name string, def float32) float32 {
	if value, set := s[name]; set {
		float, err := value.Float()
		if err != nil {
			panic(err)
		}
		return float
	}
	return def
}

func (s Style) FloatRef(name string, def float32) FloatRef {
	return FloatRef{style: s, name: name, def: def}
}

type Variable interface {
	Float() (float32, error)
	Color() (render.Color, error)
}

type ColorValue render.Color

func (c ColorValue) Float() (float32, error)      { return 0, ErrIllegalCast }
func (c ColorValue) Color() (render.Color, error) { return render.Color(c), nil }
func Color(color render.Color) Variable           { return ColorValue(color) }

type FloatValue float32

func (f FloatValue) Float() (float32, error)      { return float32(f), nil }
func (f FloatValue) Color() (render.Color, error) { return render.Black, ErrIllegalCast }
func Float(f float32) Variable                    { return FloatValue(f) }

type FloatRef struct {
	style Style
	name  string
	def   float32
}

func (f FloatRef) Val() float32 {
	return f.style.Float(f.name, f.def)
}
