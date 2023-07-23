package object

import "github.com/johanhenriksson/goworld/core/transform"

type ghost struct {
	object
	target transform.T
}

func Ghost(name string, target transform.T) Object {
	ghost := &ghost{
		object: emptyObject("Ghost:" + name),
		target: target,
	}
	ghost.transform = target
	return ghost
}

func (g *ghost) setParent(parent Object) {
	g.component.setParent(parent)
	// do not modify transform hierarchy
}

type float struct {
	object
}

// Floating objects are not part of the transform heirarchy
func Float(name string) Object {
	float := &float{
		object: emptyObject(name),
	}
	float.transform.SetParent(nil)
	return float
}

func (f *float) setParent(parent Object) {
	f.component.setParent(parent)
	// do not modify transform hierarchy
}
