package events

type Data any

type Event[T Data] struct {
	callbacks []func(T)
}

func New[T Data]() Event[T] {
	return Event[T]{}
}

func (e Event[T]) Emit(event T) {
	for _, callback := range e.callbacks {
		if callback != nil {
			callback(event)
		}
	}
}

func (e *Event[T]) Subscribe(handler func(T)) func() {
	unsubscriber := func(id int) func() {
		called := false
		return func() {
			if called {
				// its not safe to call the unsubscriber multiple times since the id might be reused
				panic("unsubscriber called multiple times")
			}
			e.callbacks[id] = nil
			called = true
		}
	}

	for id, callback := range e.callbacks {
		if callback == nil {
			e.callbacks[id] = handler
			return unsubscriber(id)
		}
	}

	e.callbacks = append(e.callbacks, handler)
	return unsubscriber(len(e.callbacks) - 1)
}
