package attachment

import (
	"github.com/vkngwrapper/core/v2/core1_0"
)

var BlendMix = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: core1_0.BlendOpAdd,
		SrcFactor: core1_0.BlendFactorSrcAlpha,
		DstFactor: core1_0.BlendFactorOneMinusSrcAlpha,
	},
	Alpha: BlendOp{
		Operation: core1_0.BlendOpAdd,
		SrcFactor: core1_0.BlendFactorOne,
		DstFactor: core1_0.BlendFactorZero,
	},
}

var BlendAdditive = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: core1_0.BlendOpAdd,
		SrcFactor: core1_0.BlendFactorOne,
		DstFactor: core1_0.BlendFactorOne,
	},
	Alpha: BlendOp{
		Operation: core1_0.BlendOpAdd,
		SrcFactor: core1_0.BlendFactorOne,
		DstFactor: core1_0.BlendFactorZero,
	},
}

var BlendMultiply = Blend{
	Enabled: true,
	Color: BlendOp{
		Operation: core1_0.BlendOpAdd,
		SrcFactor: core1_0.BlendFactorSrcAlpha,
		DstFactor: core1_0.BlendFactorOneMinusSrcAlpha,
	},
	Alpha: BlendOp{
		Operation: core1_0.BlendOpAdd,
		SrcFactor: core1_0.BlendFactorSrcAlpha,
		DstFactor: core1_0.BlendFactorOneMinusSrcAlpha,
	},
}
