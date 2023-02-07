package transform_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLabel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transform Suite")
}

var _ = Describe("", func() {
	It("initializes properly", func() {
		t := transform.Identity()
		Expect(t.Forward()).To(Equal(vec3.UnitZ))
	})
})
