package object_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/johanhenriksson/goworld/core/object"
)

type EventTester struct {
	object.Object
	OnEnableCalls  int
	OnDisableCalls int
}

var _ object.EnableHandler = &EventTester{}
var _ object.DisableHandler = &EventTester{}

func (e *EventTester) OnEnable()  { e.OnEnableCalls++ }
func (e *EventTester) OnDisable() { e.OnDisableCalls++ }

var _ = Describe("events test", func() {
	var scene object.Object
	var tester *EventTester
	pool := object.NewPool()

	BeforeEach(func() {
		scene = object.Scene(pool)
		tester = object.New(pool, "Tester", &EventTester{})
		Expect(tester.Active()).To(BeFalse())
	})

	Context("OnEnable tests", func() {
		It("activates when directly attached to a scene", func() {
			object.Attach(scene, tester)
			Expect(tester.OnEnableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeTrue())
		})

		It("does not activate when attached to a parent thats not attached to a scene", func() {
			parent := object.Empty(pool, "")
			object.Attach(parent, tester)
			Expect(tester.OnEnableCalls).To(Equal(0))
			Expect(tester.Active()).To(BeFalse())
		})

		It("recursively activates when attached to a scene", func() {
			parent := object.Empty(pool, "")
			object.Attach(parent, tester)
			object.Attach(scene, parent)
			Expect(tester.OnEnableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeTrue())
		})

		It("does not activate when attaching a disabled object to a scene", func() {
			object.Disable(tester)
			object.Attach(scene, tester)
			Expect(tester.OnEnableCalls).To(Equal(0))
			Expect(tester.Active()).To(BeFalse())
		})

		It("activates when enabling objects attached to a scene", func() {
			object.Disable(tester)
			object.Attach(scene, tester)
			object.Enable(tester)
			Expect(tester.OnEnableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeTrue())
		})
	})

	Context("OnDisable event", func() {
		It("deactivates objects when detaching", func() {
			object.Attach(scene, tester)
			object.Detach(tester)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
		})

		It("recursively deactivates children when detaching", func() {
			parent := object.Empty(pool, "")
			object.Attach(parent, tester)
			object.Attach(scene, parent)
			object.Detach(parent)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
		})

		It("deactivates object on disable", func() {
			object.Attach(scene, tester)
			object.Disable(tester)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
		})

		It("deactivates but does not disable recurisvely", func() {
			parent := object.Empty(pool, "")
			object.Attach(parent, tester)
			object.Attach(scene, parent)
			object.Disable(parent)
			Expect(tester.OnDisableCalls).To(Equal(1))
			Expect(tester.Active()).To(BeFalse())
			Expect(tester.Enabled()).To(BeTrue())
		})
	})
})
