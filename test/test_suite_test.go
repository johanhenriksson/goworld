package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/math/random"
)

func TestSuiteTest(t *testing.T) {
	random.Seed(0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}
