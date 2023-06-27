package object

type ActivateHandler interface {
	T
	OnActivate()
}

type AttachHandler interface {
	T
	OnAttach(T)
}

type DetachHandler interface {
	T
	OnDetach(T)
}
