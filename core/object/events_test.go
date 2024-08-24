package object_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/johanhenriksson/goworld/core/object"
)

type EventTester struct {
	Object
	OnEnableCalls  int
	OnDisableCalls int
}

var _ EnableHandler = &EventTester{}
var _ DisableHandler = &EventTester{}

func (e *EventTester) OnEnable()  { e.OnEnableCalls++ }
func (e *EventTester) OnDisable() { e.OnDisableCalls++ }

var _ = Describe("events test", func() {
	var scene Object
	var tester *EventTester
	pool := NewPool()

	BeforeEach(func() {
		scene = Scene(pool)
		tester = NewObject(pool, "Tester", &EventTester{})
		Expect(tester.Active()).To(BeFalse())
	})

	Context("OnEnable tests", func() {
		It("activates when directly attached to a scene", func() {
			Attach(scene, tester)
			Expect(tester.OnEnableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeTrue())
		})

		It("does not activate when attached to a parent thats not attached to a scene", func() {
			parent := Empty(pool, "")
			Attach(parent, tester)
			Expect(tester.OnEnableCalls).To(Equal(0))
			Expect(tester.Active()).To(BeFalse())
		})

		It("recursively activates when attached to a scene", func() {
			parent := Empty(pool, "")
			Attach(parent, tester)
			Attach(scene, parent)
			Expect(tester.OnEnableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeTrue())
		})

		It("does not activate when attaching a disabled object to a scene", func() {
			Disable(tester)
			Attach(scene, tester)
			Expect(tester.OnEnableCalls).To(Equal(0))
			Expect(tester.Active()).To(BeFalse())
		})

		It("activates when enabling objects attached to a scene", func() {
			Disable(tester)
			Attach(scene, tester)
			Enable(tester)
			Expect(tester.OnEnableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeTrue())
		})
	})

	Context("OnDisable event", func() {
		It("deactivates objects when detaching", func() {
			Attach(scene, tester)
			Detach(tester)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
		})

		It("recursively deactivates children when detaching", func() {
			parent := Empty(pool, "")
			Attach(parent, tester)
			Attach(scene, parent)
			Detach(parent)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
		})

		It("deactivates object on disable", func() {
			Attach(scene, tester)
			Disable(tester)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
		})

		It("deactivates but does not disable recurisvely", func() {
			parent := Empty(pool, "")
			Attach(parent, tester)
			Attach(scene, parent)
			Disable(parent)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
			Expect(tester.Enabled()).To(BeTrue())
		})
	})
})
