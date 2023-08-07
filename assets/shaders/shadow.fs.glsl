#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/lighting.glsl"

IN(0, float, depth)

void main() 
{
	// exponential depth
	gl_FragDepth = exp(SHADOW_POWER * in_depth) / exp(SHADOW_POWER);
}
