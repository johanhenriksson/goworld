package hooks

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func UseString(initial string) (string, func(string)) {
	state, setState := UseState(initial)
	setter := func(new string) { setState(new) }
	value := state.(string)
	return value, setter
}

func UseInt(initial int) (int, func(int)) {
	state, setState := UseState(initial)
	setter := func(new int) { setState(new) }
	value := state.(int)
	return value, setter
}

func UseFloat(initial float32) (float32, func(float32)) {
	state, setState := UseState(initial)
	setter := func(new float32) { setState(new) }
	value := state.(float32)
	return value, setter
}

func UseVec2(initial vec2.T) (vec2.T, func(vec2.T)) {
	state, setState := UseState(initial)
	setter := func(new vec2.T) { setState(new) }
	value := state.(vec2.T)
	return value, setter
}

func UseVec3(initial vec3.T) (vec3.T, func(vec3.T)) {
	state, setState := UseState(initial)
	setter := func(new vec3.T) { setState(new) }
	value := state.(vec3.T)
	return value, setter
}
