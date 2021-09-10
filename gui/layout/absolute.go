package layout

import "fmt"
import "github.com/johanhenriksson/goworld/math/vec2"

type Absolute struct {
	Width  float32
	Height float32
}

func (a Absolute) Flow(r Layoutable) {
	for _, item := range r.Children() {
		w := item.Width().Resolve(a.Width)
		h := item.Height().Resolve(a.Height)
		fmt.Println("Flow", item.Key(), item.Width(), item.Height(), w, h)
		item.Resize(vec2.New(w, h))
	}
}
