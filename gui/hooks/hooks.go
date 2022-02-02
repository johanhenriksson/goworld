package hooks

type State interface{}
type Setter func(State)
type Effect func()

type StateCallback func()

var hookState = []State{}
var nextHook = 0

var notify StateCallback

func SetCallback(cb StateCallback) {
	notify = cb
}

func Reset() {
	nextHook = 0
}

func UseState(initial State) (State, Setter) {
	index := nextHook
	nextHook++

	// store state
	state := initial
	if len(hookState) > index {
		state = hookState[index]
	} else {
		hookState = append(hookState, state)
	}

	setter := func(new State) {
		hookState[index] = new
		if notify != nil {
			notify()
		}
	}
	return state, setter
}

func UseEffect(callback Effect, deps ...State) {
	index := nextHook
	nextHook++

	noDeps := len(deps) == 0
	changed := false
	if len(hookState) > index {
		prev := hookState[index].([]State)
		for i, dep := range deps {
			if prev[i] != dep {
				changed = true
				break
			}
		}
	} else {
		hookState[index] = deps
	}

	if noDeps || changed {
		callback()
	}
}
