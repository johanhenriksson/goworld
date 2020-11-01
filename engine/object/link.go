package object

// Link behaviours to objects
type Link struct {
	*T
	name    string
	enabled bool
}

func NewLink(name string) *Link {
	return &Link{
		name:    name,
		enabled: true,
	}
}

func (b *Link) String() string { return b.name }

// Parent refers to the parent object
func (b *Link) Parent() *T { return b.T }

// SetParent sets the parent object
func (b *Link) SetParent(o *T) { b.T = o }

// Active returns true if the object is active
func (b *Link) Active() bool { return b.enabled }

// SetActive sets the active state of the object
func (b *Link) SetActive(active bool) { b.enabled = active }

// Update the object
func (b *Link) Update(dt float32) {}

// Collect query results on an object
func (b *Link) Collect(q *Query) {}
