package engine

import (
	"fmt"
	"runtime"

	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

func NewStatsGUI() gui.Fragment {
	lastAlloc := uint64(0)
	timer := NewFrameCounter(100)

	return gui.NewFragment(gui.FragmentArgs{
		Slot: "gui",
		Render: func() node.T {
			timer.Update()
			timings := timer.Sample()
			avgFps := 1.0 / timings.Average

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			heapAlloc := m.Alloc / 1024 / 1024
			frameAlloc := (m.TotalAlloc - lastAlloc) / 1024
			lastAlloc = m.TotalAlloc

			return rect.New("stats", rect.Props{
				Style: rect.Style{
					Position: style.Absolute{
						Bottom: style.Px(4),
						Right:  style.Px(10),
					},
					Layout:     style.Column{},
					AlignItems: style.AlignEnd,
				},
				Children: []node.T{
					label.New("fps", label.Props{
						Text: fmt.Sprintf("fps=%.1f", avgFps),
						Style: label.Style{
							Color: color.White,
						},
					}),
					label.New("mem", label.Props{
						Text: fmt.Sprintf("heap=%dmb alloc=%dkb gc=%d", heapAlloc, frameAlloc, m.NumGC),
						Style: label.Style{
							Color: color.White,
						},
					}),
				},
			})
		},
	})
}
