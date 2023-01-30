package attachment

import vk "github.com/vulkan-go/vulkan"

var BlendMix = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: vk.BlendOpAdd,
		SrcFactor: vk.BlendFactorSrcAlpha,
		DstFactor: vk.BlendFactorOneMinusSrcAlpha,
	},
	Alpha: BlendOp{
		Operation: vk.BlendOpAdd,
		SrcFactor: vk.BlendFactorOne,
		DstFactor: vk.BlendFactorZero,
	},
}

var BlendAdditive = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: vk.BlendOpAdd,
		SrcFactor: vk.BlendFactorOne,
		DstFactor: vk.BlendFactorOne,
	},
	Alpha: BlendOp{
		Operation: vk.BlendOpAdd,
		SrcFactor: vk.BlendFactorOne,
		DstFactor: vk.BlendFactorZero,
	},
}
