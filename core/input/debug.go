package input

import (
	"log"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
)

type nopHandler struct{}

func (h *nopHandler) KeyEvent(e keys.Event)    {}
func (h *nopHandler) MouseEvent(e mouse.Event) {}

func NopHandler() Handler {
	return &nopHandler{}
}

type debugger struct {
	Handler
}

func DebugMiddleware(next Handler) Handler {
	return &debugger{next}
}

func (d debugger) KeyEvent(e keys.Event) {
	log.Printf("%+v\n", e)
	if d.Handler != nil {
		d.Handler.KeyEvent(e)
	}
}

func (d debugger) MouseEvent(e mouse.Event) {
	log.Printf("%+v\n", e)
	if d.Handler != nil {
		d.Handler.MouseEvent(e)
	}
}
