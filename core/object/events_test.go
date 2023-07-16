package object_test

import (
	"github.com/johanhenriksson/goworld/core/object"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type EventTester struct {
	object.Object
	OnEnableCalled  bool
	OnDisableCalled bool
}

var _ object.EnableHandler = &EventTester{}
var _ object.DisableHandler = &EventTester{}

func (e *EventTester) OnEnable()  { e.OnEnableCalled = true }
func (e *EventTester) OnDisable() { e.OnDisableCalled = true }

var _ = Describe("events test", func() {
	var scene object.Object
	var tester *EventTester

	BeforeEach(func() {
		scene = object.Scene()
		tester = object.New("Tester", &EventTester{})
		Expect(tester.Active()).To(BeFalse())
	})

	Context("OnEnable tests", func() {
		It("activates when directly attached to a scene", func() {
			object.Attach(scene, tester)
			Expect(tester.OnEnableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeTrue())
		})

		It("does not activate when attached to a parent thats not attached to a scene", func() {
			parent := object.Empty("")
			object.Attach(parent, tester)
			Expect(tester.OnEnableCalled).To(BeFalse())
			Expect(tester.Active()).To(BeFalse())
		})

		It("recursively activates when attached to a scene", func() {
			parent := object.Empty("")
			object.Attach(parent, tester)
			object.Attach(scene, parent)
			Expect(tester.OnEnableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeTrue())
		})

		It("does not activate when attaching a disabled object to a scene", func() {
			object.Disable(tester)
			object.Attach(scene, tester)
			Expect(tester.OnEnableCalled).To(BeFalse())
			Expect(tester.Active()).To(BeFalse())
		})

		It("activates when enabling objects attached to a scene", func() {
			object.Disable(tester)
			object.Attach(scene, tester)
			object.Enable(tester)
			Expect(tester.OnEnableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeTrue())
		})
	})

	Context("OnDisable event", func() {
		It("deactivates objects when detaching", func() {
			object.Attach(scene, tester)
			object.Detach(tester)
			Expect(tester.OnDisableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeFalse())
		})

		It("recursively deactivates children when detaching", func() {
			parent := object.Empty("")
			object.Attach(parent, tester)
			object.Attach(scene, parent)
			object.Detach(parent)
			Expect(tester.OnDisableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeFalse())
		})

		It("deactivates object on disable", func() {
			object.Attach(scene, tester)
			object.Disable(tester)
			Expect(tester.OnDisableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeFalse())
		})

		It("deactivates but does not disable recurisvely", func() {
			parent := object.Empty("")
			object.Attach(parent, tester)
			object.Attach(scene, parent)
			object.Disable(parent)
			Expect(tester.OnDisableCalled).To(BeTrue())
			Expect(tester.Active()).To(BeFalse())
			Expect(tester.Enabled()).To(BeTrue())
		})
	})
})
