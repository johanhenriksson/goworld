package engine

type Element interface {
}

type Value interface {
	String() string
	Int() int
}

type String struct {
	Value string
}

func (s *String) String() (string, error) {
	return s.Value, nil
}

type Rect struct {
	Layout   Value
	Align    Value
	Visible  Value
	Children []Element
}

type State map[string]Value

func render(root Element, state State) {

}
