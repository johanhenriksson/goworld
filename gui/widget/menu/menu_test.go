package menu_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget/menu"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestButton(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Menu Widget Suite")
}

var _ = Describe("Menu widget", func() {
	It("renders", func() {
		render := func() node.T {
			return rect.New("gui", rect.Props{
				Children: []node.T{
					makeMenu(),
				},
			})
		}

		tree := render()
		Expect(tree.Children()).To(HaveLen(1))
	})
})

func makeMenu() node.T {
	return menu.Menu("gui-menu", menu.Props{
		Style: menu.Style{
			Color:      color.RGB(0.76, 0.76, 0.76),
			HoverColor: color.RGB(0.85, 0.85, 0.85),
			TextColor:  color.Black,
		},

		Items: []menu.ItemProps{
			{
				Key:   "menu-file",
				Title: "File",
				Items: []menu.ItemProps{
					{
						Key:   "file-exit",
						Title: "Exit",
					},
				},
			},
		},
	})
}
