package render

// Pass is a type identifier for a render draw pass.
type Pass string

// Various draw pass constants
const (
	Geometry    Pass = "geometry"
	Shadow           = "shadow"
	Line             = "lines"
	Particles        = "particles"
	Forward          = "forward"
	Lights           = "light"
	Postprocess      = "postprocess"
	UI               = "ui"
)

type Passes []Pass

func (p Passes) Includes(other Pass) bool {
	for _, pass := range p {
		if pass == other {
			return true
		}
	}
	return false
}

func (p *Passes) Add(passes ...Pass) {
	for _, pass := range passes {
		if p.Includes(pass) {
			return
		}
		*p = append(*p, pass)
	}
}

func (p *Passes) Set(passes ...Pass) {
	*p = passes
}
