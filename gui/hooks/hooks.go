package hooks

type Setter func(any)
type Effect func()

type StateCallback func()

var hookState = []any{}
var nextHook = 0

type State struct {
	data []any
	next int
}

var state *State = nil

func Enable(new *State) {
	state = new
	state.next = 0
}

func Disable() {
	state = nil
}

func UseState[T any](initial T) (T, func(T)) {
	if state == nil {
		panic("no active hook state")
	}

	index := state.next
	state.next++

	// store state
	value := initial
	if len(hookState) > index {
		value = state.data[index].(T)
	} else {
		state.data = append(state.data, value)
	}

	setter := func(new T) {
		state.data[index] = new
	}
	return value, setter
}

func UseEffect(callback Effect, deps ...any) {
	if state == nil {
		panic("no active hook state")
	}

	index := state.next
	state.next++

	noDeps := len(deps) == 0
	changed := false
	if len(state.data) > index {
		prev := state.data[index].([]any)
		for i, dep := range deps {
			if prev[i] != dep {
				changed = true
				break
			}
		}
	} else {
		state.data[index] = deps
	}

	if noDeps || changed {
		callback()
	}
}
