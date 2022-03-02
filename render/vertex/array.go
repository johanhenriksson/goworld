package vertex

type Array interface {
	Bind() error
	Unbind()
	Delete()
	Indexed() bool
	Draw() error

	SetIndexSize(int)
	SetElements(int)
	SetPointers(Pointers)

	Buffer(name string, data any)
}
