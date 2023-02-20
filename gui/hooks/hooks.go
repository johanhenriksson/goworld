package hooks

type Effect func()

func UseState[T any](initial T) (T, func(T)) {
	// keep a reference to the current state
	state := getState()
	index := state.Next()

	// store state
	value := initial
	if len(state.data) > index {
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
	// keep a reference to the current state
	state := getState()
	index := state.Next()

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
		state.data = append(state.data, deps)
	}

	if noDeps || changed {
		callback()
	}
}
