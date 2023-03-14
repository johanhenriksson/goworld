package label

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/input/keys"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLabel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Label Suite")
}

var _ = Describe("", func() {
	Context("key events", func() {
		var label T
		var props Props

		BeforeEach(func() {
			props = Props{
				OnChange: func(new string) {
					props.Text = new
				},
			}
			label = new("label", props).(T)
		})

		It("adds characters", func() {
			ev1 := keys.NewCharEvent('a', keys.NoMod)
			label.KeyEvent(ev1)
			label.KeyEvent(ev1)

			label.Update(props)

			Expect(label.Cursor()).To(Equal(2))
			Expect(label.Text()).To(Equal("aa"))
		})

		It("removes on backspace", func() {
			props.Text = "ok"
			label.Update(props)
			Expect(label.Cursor()).To(Equal(2))

			label.KeyEvent(keys.NewPressEvent(keys.Backspace, keys.Press, keys.NoMod))
			label.KeyEvent(keys.NewPressEvent(keys.Backspace, keys.Press, keys.NoMod))
			label.KeyEvent(keys.NewPressEvent(keys.Backspace, keys.Press, keys.NoMod))
			label.Update(props)

			Expect(label.Text()).To(Equal(""))
			Expect(label.Cursor()).To(Equal(0))
		})

		It("removes on forward delete", func() {
			props.Text = "ok!"
			label.Update(props)

			label.KeyEvent(keys.NewPressEvent(keys.LeftArrow, keys.Press, keys.NoMod))
			label.KeyEvent(keys.NewPressEvent(keys.Delete, keys.Press, keys.NoMod))
			label.Update(props)

			Expect(label.Text()).To(Equal("ok"))
			Expect(label.Cursor()).To(Equal(2))
		})
	})
})
