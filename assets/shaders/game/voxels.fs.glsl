#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

void main() 
{
	out_diffuse = vec4(in_color.rgb * in_color.a, 1);
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
