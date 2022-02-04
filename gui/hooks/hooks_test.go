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
	output, click := SomeComponent()
	if output != "hello!" {
		t.Error("unexpected return value")
	}

	// prepare for next render
	hooks.Reset()

	click()

	output, _ = SomeComponent()
	if output != "clicked" {
		t.Error("expected state to be updated")
	}
}
