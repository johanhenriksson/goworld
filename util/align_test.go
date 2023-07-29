package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/johanhenriksson/goworld/util"
)

type AlignCase struct {
	Offset    int
	Alignment int
	Expected  int
}

var _ = Describe("align utils", func() {
	It("returns the expected alignment", func() {
		cases := []AlignCase{
			{23, 64, 64},
			{64, 64, 64},
			{72, 64, 128},
		}
		for _, testcase := range cases {
			actual := Align(testcase.Offset, testcase.Alignment)
			Expect(actual).To(Equal(testcase.Expected))
		}
	})

	It("returns errors for misaligned structs", func() {
		type FailingStruct struct {
			A bool
			B int
		}
		err := ValidateAlignment(FailingStruct{})
		Expect(err).To(HaveOccurred())
	})

	It("validates aligned structs", func() {
		type PassingStruct struct {
			A int
			B float32
		}
		err := ValidateAlignment(PassingStruct{})
		Expect(err).ToNot(HaveOccurred())
	})
})
