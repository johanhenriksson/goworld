#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

OBJECT(1, object)

void main() 
{
	out_diffuse = in_color;
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
