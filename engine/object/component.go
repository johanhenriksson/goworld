package object

type Component interface {
	String() string
	Update(float32)

	Active() bool
	SetActive(bool)

	Parent() *T
	SetParent(*T)
	Collect(*Query)
}
