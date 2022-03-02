package vertex

type Buffer interface {
	Bind()
	Unbind()
	Delete()
	Buffer(data any)
}
