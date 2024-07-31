package texture

import "github.com/vkngwrapper/core/v2/core1_0"

type Filter core1_0.Filter

const FilterNearest = Filter(core1_0.FilterNearest)
const FilterLinear = Filter(core1_0.FilterLinear)

type Wrap core1_0.SamplerAddressMode

const WrapClamp = Wrap(core1_0.SamplerAddressModeClampToEdge)
const WrapRepeat = Wrap(core1_0.SamplerAddressModeRepeat)
const WrapMirror = Wrap(core1_0.SamplerAddressModeMirroredRepeat)
const WrapBorder = Wrap(core1_0.SamplerAddressModeClampToBorder)
