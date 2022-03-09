package descriptor_test

import (
	"testing"

	. "github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"

	vk "github.com/vulkan-go/vulkan"
)

type TestSet struct {
	Set
	Diffuse *Sampler
}

type MockSet struct{}

var _ Set = &MockSet{}

func (m *MockSet) Ptr() vk.DescriptorSet { return nil }

func (m *MockSet) Write(w vk.WriteDescriptorSet) {}

func TestParseDescriptors(t *testing.T) {
	set := TestSet{
		Diffuse: &Sampler{
			Binding: 0,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	}
	desc := ParseDescriptors(&set)
	if _, ok := desc["diffuse"]; !ok {
		t.Error("expected to find diffuse descriptor")
	}
}

func TestCopySet(t *testing.T) {
	desc := TestSet{
		Diffuse: &Sampler{
			Binding: 1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	}
	set := &MockSet{}

	desc2 := BindSet(&desc, set)
	if desc2 == nil {
		t.Error("set copy should not be nil")
	}
	if desc2.Diffuse == nil {
		t.Error("expected diffuse to be set on the copy")
	}
	if desc2.Diffuse == desc.Diffuse {
		t.Error("expected a new copy of diffuse")
	}
	if desc2.Diffuse.Binding != 1 {
		t.Error("expected diffuse data to be copied")
	}
	if desc2.Set != set {
		t.Error("expected Set to be MockSet, was", desc2.Set)
	}
}
