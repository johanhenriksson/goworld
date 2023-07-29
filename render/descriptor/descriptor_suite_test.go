package descriptor_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDescriptor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "render/descriptor")
}
