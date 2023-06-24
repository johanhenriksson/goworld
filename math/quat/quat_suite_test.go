package quat_test

import (
	"testing"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestQuat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Quat Suite")
}

const e = 0.001

var _ = Describe("quaternion", func() {
	Context("euler angles", func() {
		It("converts back and forth", func() {
			x, y, z := float32(10), float32(20), float32(30)
			q := quat.Euler(x, y, z)
			r := q.Euler()
			GinkgoWriter.Println(x, y, z, r)
			Expect(r.X).To(BeNumerically("~", x, e), "wrong x rotation")
			Expect(r.Y).To(BeNumerically("~", y, e), "wrong y rotation")
			Expect(r.Z).To(BeNumerically("~", z, e), "wrong z rotation")
		})

		It("returns the expected rotation matrices", func() {
			compareMatrices := func(rotX, rotY, rotZ float32) {
				q := quat.Euler(rotX, rotY, rotZ)
				m1 := q.Mat4()
				m2 := mat4.Rotate(vec3.New(rotX, rotY, rotZ))

				GinkgoWriter.Println(rotX, rotY, rotZ)

				Expect(m1.Right().X).To(BeNumerically("~", m2.Right().X, e), "wrong right vector (x)")
				Expect(m1.Right().Y).To(BeNumerically("~", m2.Right().Y, e), "wrong right vector (y)")
				Expect(m1.Right().Z).To(BeNumerically("~", m2.Right().Z, e), "wrong right vector (z)")

				Expect(m1.Up().X).To(BeNumerically("~", m2.Up().X, e), "wrong up vector (x)")
				Expect(m1.Up().Y).To(BeNumerically("~", m2.Up().Y, e), "wrong up vector (y)")
				Expect(m1.Up().Z).To(BeNumerically("~", m2.Up().Z, e), "wrong up vector (z)")

				Expect(m1.Forward().X).To(BeNumerically("~", m2.Forward().X, e), "wrong forward vector (x)")
				Expect(m1.Forward().Y).To(BeNumerically("~", m2.Forward().Y, e), "wrong forward vector (y)")
				Expect(m1.Forward().Z).To(BeNumerically("~", m2.Forward().Z, e), "wrong forward vector (z)")
			}
			compareMatrices(0, 0, 0)
			compareMatrices(10, 20, 30)
			compareMatrices(170, 60, 99)
			compareMatrices(280, 280, 280)
		})
	})
})
