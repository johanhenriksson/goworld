package device

import "github.com/vkngwrapper/core/v2/core1_0"

type Queue struct {
	ptr    core1_0.Queue
	family int
}

func (q Queue) Ptr() core1_0.Queue {
	return q.ptr
}

func (q Queue) FamilyIndex() int {
	return q.family
}

func (q Queue) Index() int {
	return 0
}
