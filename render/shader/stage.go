package shader

import "github.com/vkngwrapper/core/v2/core1_0"

type ShaderStage core1_0.ShaderStageFlags

const (
	StageAll      = ShaderStage(core1_0.StageAll)
	StageVertex   = ShaderStage(core1_0.StageVertex)
	StageFragment = ShaderStage(core1_0.StageFragment)
	StageCompute  = ShaderStage(core1_0.StageCompute)
)

func (s ShaderStage) String() string {
	return s.flags().String()
}

// flags returns the Vulkan-native representation
func (s ShaderStage) flags() core1_0.ShaderStageFlags {
	return core1_0.ShaderStageFlags(s)
}
