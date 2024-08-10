#version 450

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

OBJECT(1, object)
SAMPLER_ARRAY(3, textures)

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = object.textures[0];

	out_diffuse = vec4(texture_array(textures, texture0, texcoord0).rgb, 1);
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
