package object

// Link behaviours to objects
type Link struct {
	*T
}

// NewLink creates a new component-object link
func NewLink(parent *T) *Link {
	return &Link{parent}
}

// Parent refers to the parent object
func (b *Link) Parent() *T { return b.T }

// SetParent sets the parent object
func (b *Link) SetParent(o *T) Component {
	b.T = o
	return b
}

// Update the object
func (b *Link) Update(dt float32) {}

// Collect query results on an object
func (b *Link) Collect(q *Query) {}
