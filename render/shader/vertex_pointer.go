package shader

type Pointer interface {
	String() string
	Enable()
	Disable()
}
