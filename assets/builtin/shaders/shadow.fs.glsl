#version 450

#include "lib/common.glsl"
#include "lib/lighting.glsl"

IN(0, float, depth)

void main() 
{
	// exponential depth
	gl_FragDepth = exp(SHADOW_POWER * in_depth) / exp(SHADOW_POWER);
}
