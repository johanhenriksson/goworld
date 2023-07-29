package descriptor_test

import (
	. "github.com/johanhenriksson/goworld/render/descriptor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vkngwrapper/core/v2/core1_0"
)

var _ = Describe("descriptor struct parsing", func() {
	type TestSet struct {
		Set
		Diffuse *Sampler
	}

	It("correctly parses valid structs", func() {
		set := TestSet{
			Diffuse: &Sampler{
				Stages: core1_0.StageAll,
			},
		}
		desc, err := ParseDescriptorStruct(&set)
		Expect(err).ToNot(HaveOccurred())
		Expect(desc).To(HaveLen(1), "expected to find diffuse descriptor")
	})

	It("rejects unset descriptor fields", func() {
		set := TestSet{
			Diffuse: nil,
		}
		_, err := ParseDescriptorStruct(&set)
		Expect(err).To(HaveOccurred())
	})

	It("rejects non-pointer fields", func() {
		type FailSet struct {
			Set
			Diffuse Sampler
		}
		set := FailSet{}
		_, err := ParseDescriptorStruct(&set)
		Expect(err).To(HaveOccurred())
	})
})
