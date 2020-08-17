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
	_, dh := rect.DesiredSize(300, 300)

	if dh != 32.0 {
		t.Errorf("expected height %f", dh)
	}
}
