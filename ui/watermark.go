package ui

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Watermark struct {
	*Text
	wnd *engine.Window
}

func NewWatermark(wnd *engine.Window) *Watermark {
	w := &Watermark{
		Text: NewText("", Style{}),
		wnd:  wnd,
	}
	w.SetPosition(vec2.New(10, float32(wnd.Height-30)))
	w.SetZIndex(1000)
	return w
}

func (w *Watermark) Draw(args engine.DrawArgs) {
	w.Set(fmt.Sprintf("goworld | %s | %.0f fps", time.Now().Format("2006-01-02 15:04"), w.wnd.FPS))
	w.Text.Draw(args)
}
