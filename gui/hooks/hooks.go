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

	// if we have no dependencies, run the callback every frame
	noDeps := len(deps) == 0
	if noDeps {
		callback()
		return
	}

	// check if any dependencies have changed
	changed := false
	if len(state.data) > index {
		// previous state exists - compare it
		prev := state.data[index].([]any)
		for i, dep := range deps {
			if prev[i] != dep {
				changed = true
				// update previous value so that the callback wont run again next time
				prev[i] = dep
			}
		}
	} else {
		// no previous state exists yet
		state.data = append(state.data, deps)
	}

	// finally, run the callback if anything changed
	if changed {
		callback()
	}
}
