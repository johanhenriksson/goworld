package object

type ActivateHandler interface {
	Component
	OnActivate()
}

type DeactivateHandler interface {
	Component
	OnDeactivate()
}

type AttachHandler interface {
	Object
	// OnAttach is called when a component or child is attached to the group
	OnAttach(Component)
}

type DetachHandler interface {
	Object
	// OnDetach is called when a component or child is detached from the group
	OnDetach(Component)
}

type AttachedHandler interface {
	Component
	// OnAttached is called when the component or group is attached to a parent group
	OnAttached(Object)
}

type DetachedHandler interface {
	Component
	// OnDetached is called when the component or group is detached from a parent group
	OnDetached(Object)
}

type ChildEventHandler interface {
	Component
	OnAddChild(Component)
	OnRemoveChild(Component)
}

type ParentEventHandler interface {
	Component
	OnAttachTo(Object)
	OnDetachFrom(Object)
}

type SiblingEventHandler interface {
	Component
	OnAddSibling(Component)
	OnRemoveSibling(Component)
}
