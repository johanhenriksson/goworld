package vertex

import "github.com/vkngwrapper/core/v2/core1_0"

const (
	CullNone  = CullMode(core1_0.CullModeFlags(0))
	CullFront = CullMode(core1_0.CullModeFront)
	CullBack  = CullMode(core1_0.CullModeBack)
)

type CullMode core1_0.CullModeFlags

func (c CullMode) flags() core1_0.CullModeFlags {
	return core1_0.CullModeFlags(c)
}

func (c CullMode) String() string {
	return c.flags().String()
}
