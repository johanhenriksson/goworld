package object

func Attach(parent, child T) {
	Detach(child)
	child.setParent(parent)
	parent.attach(child)
}

func Detach(child T) {
	if child.Parent() == nil {
		return
	}
	child.Parent().detach(child)
	child.setParent(nil)
}

func Root(obj T) T {
	for obj.Parent() != nil {
		obj = obj.Parent()
	}
	return obj
}

func GetInParents[K T](root T) (K, bool) {
	if k, ok := root.(K); ok {
		return k, true
	}
	if root.Parent() != nil {
		return GetInParents[K](root.Parent())
	}
	var empty K
	return empty, false
}
