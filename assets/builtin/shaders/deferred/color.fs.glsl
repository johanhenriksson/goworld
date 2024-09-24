#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"
#include "lib/deferred_fragment.glsl"

OBJECT(1, object)

void main() 
{
	out_diffuse = in_color;
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
