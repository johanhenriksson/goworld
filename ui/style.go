package ui

import (
	"fmt"
	"github.com/johanhenriksson/goworld/render"
)

type Styled interface {
	Float(string, float32) float32
	String(string, string) string
	Color(string, render.Color) render.Color
}

// ErrIllegalCast occurs when attempting to cast a style variable to an incompatible type
var ErrIllegalCast = fmt.Errorf("illegal cast")

// NoStyle is the empty style sheet
var NoStyle = Style{}

// Style holds UI styling information
type Style map[string]Variable

func (s Style) Extend(s2 Style) Style {
	out := Style{}
	for k, v := range s {
		out[k] = v
	}
	for k, v := range s2 {
		out[k] = v
	}
	return out
}

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

// String returns a string value from the styles
func (s Style) String(name string, def string) string {
	if value, set := s[name]; set {
		float, err := value.String()
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
	String() (string, error)
}

type ColorValue render.Color

func (c ColorValue) Float() (float32, error)      { return 0, ErrIllegalCast }
func (c ColorValue) Color() (render.Color, error) { return render.Color(c), nil }
func (c ColorValue) String() (string, error)      { return render.Color(c).String(), nil }
func Color(color render.Color) Variable           { return ColorValue(color) }

type FloatValue float32

func (f FloatValue) Float() (float32, error)      { return float32(f), nil }
func (f FloatValue) Color() (render.Color, error) { return render.Black, ErrIllegalCast }
func (f FloatValue) String() (string, error)      { return fmt.Sprintf("%f", f), nil }
func Float(f float32) Variable                    { return FloatValue(f) }

type StringValue string

func (s StringValue) Float() (float32, error)      { return 0, ErrIllegalCast }
func (s StringValue) Color() (render.Color, error) { return render.Black, ErrIllegalCast }
func (s StringValue) String() (string, error)      { return string(s), nil }
func String(str string) Variable                   { return StringValue(str) }

type FloatRef struct {
	style Style
	name  string
	def   float32
}

func (f FloatRef) Val() float32 {
	return f.style.Float(f.name, f.def)
}
