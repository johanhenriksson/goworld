#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

layout (location = 0) in float depth;

layout (location = 0) out vec4 fragColor;

void main() 
{
	fragColor = vec4(0);

	// exponential depth
	gl_FragDepth = exp(SHADOW_POWER * depth) / exp(SHADOW_POWER);
}
