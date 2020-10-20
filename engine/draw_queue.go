package engine

type DrawQueue struct {
	items []DrawCommand
}

type DrawCommand struct {
	Component Component
	Args      DrawArgs
}

func NewDrawQueue() *DrawQueue {
	capacity := 1024
	return &DrawQueue{
		items: make([]DrawCommand, 0, capacity),
	}
}

func (q *DrawQueue) Add(component Component, args DrawArgs) {
	q.items = append(q.items, DrawCommand{component, args})
}

func (q *DrawQueue) Clear() {
	// clear slice, but keep the memory
	q.items = q.items[:0]
}
