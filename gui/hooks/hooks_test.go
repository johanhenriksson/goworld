package hooks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/gui/hooks"
)

func TestHooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gui/hooks")
}

func SomeComponent() (string, func()) {
	title, setTitle := hooks.UseState("hello!")
	click := func() {
		setTitle("clicked")
	}
	return title, click
}

var _ = Describe("hooks", func() {
	It("updates and maintains state", func() {
		state := hooks.State{}
		hooks.Enable(&state)
		output, click := SomeComponent()
		hooks.Disable()
		Expect(output).To(Equal("hello!"))

		click()

		hooks.Enable(&state)
		output, _ = SomeComponent()
		hooks.Disable()
		Expect(output).To(Equal("clicked"))
	})
})
