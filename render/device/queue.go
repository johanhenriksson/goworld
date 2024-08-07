package device

import (
	"fmt"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Queue struct {
	ptr    core1_0.Queue
	flags  core1_0.QueueFlags
	family int
	index  int
}

func (q Queue) Ptr() core1_0.Queue {
	return q.ptr
}

func (q Queue) FamilyIndex() int {
	return q.family
}

func (q Queue) Index() int {
	return q.index
}

func (q Queue) Matches(flags core1_0.QueueFlags) bool {
	return q.flags&flags == flags
}

func (q Queue) String() string {
	return fmt.Sprintf("Queue{ptr: %v, flags: %v, family: %d, index: %d}", q.ptr, q.flags, q.family, q.index)
}
