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

// Write hook state. For debug/testing purposes only.
func (s *State) Write(index int, data any) {
	if len(s.data) < index+1 {
		s.data = append(s.data, make([]any, index-len(s.data)+1)...)
	}
	s.data[index] = data
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
