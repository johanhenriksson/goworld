package object

// Links components to their parent objects.
type Link struct {
	T
}

// NewLink creates a new component-object link
func NewLink(parent T) *Link {
	return &Link{parent}
}

// Parent refers to the parent object
func (b *Link) Parent() T { return b.T }

// SetParent sets the parent object
func (b *Link) SetParent(o T) {
	b.T = o
}

// Update the object. No-op when called on Links
func (b *Link) Update(dt float32) {
	// propagating the update to the linked object causes infinite recursion
	// ... do nothing ...
}

// UpdateTransform updates the transformation matrix. No-op when called on Links
func (b *Link) UpdateTransform() {
	// propagating the update to the linked object causes infinite recursion
	// ... do nothing ...
}

// Collect query results on an object. No-op when called on Links
func (b *Link) Collect(q *Query) {
	// propagating the update to the linked object causes infinite recursion
	// ... do nothing ...

	// there seems to be a pattern here !
	// perhaps indicating that the component/link design is a bad idea
}
