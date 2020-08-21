package ui

import (
	"os"
	"testing"
)

func TestFlowRect(t *testing.T) {
	os.Chdir("..")
	rect := NewRect(NoStyle)
	rect.Attach(NewText("Hello", NoStyle))
	rect.Attach(NewText("Please", NoStyle))
	desired := rect.Flow(Size{300, 300})

	if desired.Height != 32.0 {
		t.Errorf("expected height %f", desired.Height)
	}
}
