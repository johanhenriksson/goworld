package renderpass

import vk "github.com/vulkan-go/vulkan"

var BlendMultiply = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: vk.BlendOpMultiply,
		SrcFactor: vk.BlendFactorOne,
		DstFactor: vk.BlendFactorOneMinusSrcAlpha,
	},
}

var BlendAdditive = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: vk.BlendOpAdd,
		SrcFactor: vk.BlendFactorOne,
		DstFactor: vk.BlendFactorOne,
	},
}
