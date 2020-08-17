package ui

import (
	"os"
	"testing"
)

func TestFlowRect(t *testing.T) {
	os.Chdir("..")
	rect := NewRect(0, 0, 100, 50, NoStyle)
	rect.Attach(NewText("Hello", 0, 0, NoStyle))
	rect.Attach(NewText("Please", 0, 0, NoStyle))
	_, dh := rect.DesiredSize(300, 300)

	if dh != 32.0 {
		t.Errorf("expected height %f", dh)
	}
}
