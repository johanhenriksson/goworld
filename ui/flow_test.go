package ui

import (
	"os"
	"testing"

	"github.com/johanhenriksson/goworld/math/vec2"
)

func TestFlowRect(t *testing.T) {
	os.Chdir("..")
	rect := NewRect(NoStyle)
	rect.Attach(NewText("Hello", NoStyle))
	rect.Attach(NewText("Please", NoStyle))
	desired := rect.Flow(vec2.New(300, 300))

	if desired.Y != 32.0 {
		t.Errorf("expected height %f", desired.Y)
	}
}
