package events

type Data any

type Handler[T Data] func(T)

type Event[T Data] struct {
	callbacks map[any]Handler[T]
}

func New[T Data]() *Event[T] {
	return &Event[T]{
		callbacks: make(map[any]Handler[T]),
	}
}

func (e *Event[T]) Emit(event T) {
	for _, callback := range e.callbacks {
		callback(event)
	}
}

func (e *Event[T]) Subscribe(subscriber any, handler Handler[T]) {
	e.callbacks[subscriber] = handler
}

func (e *Event[T]) Unsubscribe(subscriber any) {
	delete(e.callbacks, subscriber)
}
