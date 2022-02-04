package hooks_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/hooks"
)

func SomeComponent() (string, func()) {
	title, setTitle := hooks.UseState("hello!")
	click := func() {
		setTitle("clicked")
	}
	return title, click
}

func TestHooks(t *testing.T) {
	state := hooks.State{}
	hooks.Enable(&state)
	output, click := SomeComponent()
	hooks.Disable()
	if output != "hello!" {
		t.Error("unexpected return value")
	}

	click()

	hooks.Enable(&state)
	output, _ = SomeComponent()
	hooks.Disable()
	if output != "clicked" {
		t.Error("expected state to be updated")
	}
}
