package button_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/math/vec2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestButton(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Button Suite")
}

var _ = Describe("", func() {
	renderer := node.NewRenderer(func() node.T {
		return button.New("button", button.Props{})

	})
	tree := renderer.Render(vec2.One)

	children := tree.Children()
	for _, child := range children {
		child.Children()
	}
})
