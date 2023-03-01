package font

import (
	"image"
)

type ref struct {
	key     string
	version int
	font    T
	text    string
	args    Args
}

func Ref(key string, version int, font T, text string, args Args) *ref {
	return &ref{
		key:     key,
		version: version,
		font:    font,
		text:    text,
		args:    args,
	}
}

func (r *ref) Key() string  { return r.key }
func (r *ref) Version() int { return r.version }

func (r *ref) Load() *image.RGBA {
	return r.font.Render(r.text, r.args)
}
