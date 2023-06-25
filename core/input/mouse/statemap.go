package mouse

type State interface {
	Handler

	Down(Button) bool
	Up(Button) bool
}

type state map[Button]bool

func NewState() State {
	return state{}
}

func (s state) MouseEvent(e Event) {
	if e.Action() == Press {
		s[e.Button()] = true
	}
	if e.Action() == Release {
		s[e.Button()] = false
	}
}

func (s state) Down(key Button) bool {
	if state, stored := s[key]; stored {
		return state
	}
	return false
}

func (s state) Up(key Button) bool {
	return !s.Down(key)
}
