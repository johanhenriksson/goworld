package object

type EnableHandler interface {
	Component
	OnEnable()
}

type DisableHandler interface {
	Component
	OnDisable()
}
