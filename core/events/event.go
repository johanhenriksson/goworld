package events

type Data any

type Handler[T Data] func(T)

type Event[T Data] struct {
	callbacks []Handler[T]
}

func New[T Data]() *Event[T] {
	return &Event[T]{}
}

func (e *Event[T]) Emit(event T) {
	for _, callback := range e.callbacks {
		if callback != nil {
			callback(event)
		}
	}
}

func (e *Event[T]) Subscribe(handler Handler[T]) func() {
	id := len(e.callbacks)
	e.callbacks = append(e.callbacks, handler)
	return func() {
		// unsub
		e.callbacks[id] = nil
	}
}
