package layout

import "github.com/johanhenriksson/goworld/math/vec2"

type Absolute struct{}

func (a Absolute) Flow(r Layoutable) {
	bounds := r.Size()
	for _, item := range r.Children() {
		w := item.Width().Resolve(bounds.X)
		h := item.Height().Resolve(bounds.Y)
		item.Resize(vec2.New(w, h))
	}
}
