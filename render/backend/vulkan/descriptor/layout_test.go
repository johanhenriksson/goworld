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

func TestParseDescriptors(t *testing.T) {
	set := TestSet{
		Diffuse: &Sampler{
			Binding: 0,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	}
	desc := ParseDescriptors(&set)
	if len(desc) != 1 {
		t.Error("expected to find diffuse descriptor")
	}
}
