package quat_test

import (
	. "github.com/johanhenriksson/goworld/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func TestQuat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "math/quat")
}

var _ = Describe("quaternion", func() {
	Context("euler angles", func() {
		It("converts back and forth", func() {
			x, y, z := float32(10), float32(20), float32(30)
			q := quat.Euler(x, y, z)
			r := q.Euler()
			GinkgoWriter.Println(x, y, z, r)
			Expect(r).To(ApproxVec3(vec3.New(x, y, z)), "wrong rotation")
		})
	})
})
