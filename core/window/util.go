package window

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/johanhenriksson/goworld/math"
)

func GetCurrentMonitor(window *glfw.Window) *glfw.Monitor {
	// translated to Go from https://stackoverflow.com/a/31526753
	wx, wy := window.GetPos()
	ww, wh := window.GetSize()

	bestoverlap := 0
	var bestmonitor *glfw.Monitor
	for _, monitor := range glfw.GetMonitors() {
		mode := monitor.GetVideoMode()
		mx, my := monitor.GetPos()
		mw, mh := mode.Width, mode.Height

		overlap := math.Max(0, math.Min(wx+ww, mx+mw)-math.Max(wx, mx)) *
			math.Max(0, math.Min(wy+wh, my+mh)-math.Max(wy, my))

		if bestoverlap < overlap {
			bestoverlap = overlap
			bestmonitor = monitor
		}
	}

	return bestmonitor
}
