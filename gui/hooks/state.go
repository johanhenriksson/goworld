package hooks

type State struct {
	data []any
	next int
}

func (s *State) Next() int {
	id := s.next
	s.next++
	return id
}

func Enable(new *State) {
	active = new
	active.next = 0
}

func Disable() {
	active = nil
}

var active *State = nil

func getState() *State {
	if active == nil {
		panic("no active hook state")
	}
	return active
}
