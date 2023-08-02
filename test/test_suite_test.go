package test

import (
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/johanhenriksson/goworld/math/random"
)

func TestSuiteTest(t *testing.T) {
	random.Seed(0)
	log.SetOutput(GinkgoWriter)

	// todo: clean up old failure/actual images

	RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}
