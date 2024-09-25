#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

// Varying
IN(0, flat uint, object)
IN(1, vec3, normal)
IN(2, vec3, position)
IN(3, vec4, color)

// Return Output
OUT(0, vec4, diffuse)
OUT(1, vec4, normal)
OUT(2, vec4, position)

OBJECT(1, object, in_object)
SAMPLER_ARRAY(2, textures)

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = object.textures[0];

	out_diffuse = vec4(texture_array(textures, texture0, texcoord0).rgb, 1);
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
