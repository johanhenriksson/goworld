package keys

type State interface {
	Handler

	Down(Code) bool
	Up(Code) bool

	Shift() bool
	Ctrl() bool
	Alt() bool
	Super() bool
}

type state map[Code]bool

func NewState() State {
	return state{}
}

func (s state) KeyEvent(e Event) {
	if e.Action() == Press {
		s[e.Code()] = true
	}
	if e.Action() == Release {
		s[e.Code()] = false
	}
}

func (s state) Down(key Code) bool {
	if state, stored := s[key]; stored {
		return state
	}
	return false
}

func (s state) Up(key Code) bool {
	return !s.Down(key)
}

func (s state) Shift() bool {
	return s.Down(LeftShift) || s.Down(RightShift)
}

func (s state) Alt() bool {
	return s.Down(LeftAlt) || s.Down(RightAlt)
}

func (s state) Ctrl() bool {
	return s.Down(LeftControl) || s.Down(RightControl)
}

func (s state) Super() bool {
	return s.Down(LeftSuper) || s.Down(RightSuper)
}
