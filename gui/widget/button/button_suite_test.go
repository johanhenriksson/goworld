package button_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/math/vec2"
)

func TestButton(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gui/widget/button")
}

var _ = Describe("button", func() {
	It("renders", func() {
		renderer := node.NewRenderer("test", func() node.T {
			return button.New("button", button.Props{})
		})
		tree := renderer.Render(vec2.One)

		children := tree.Children()
		for _, child := range children {
			child.Children()
		}

		// todo: incomplete test
	})
})
