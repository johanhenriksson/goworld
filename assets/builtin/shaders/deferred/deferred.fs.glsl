#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"

// Varying
IN(0, flat uint, object)
IN(1, vec3, position)
IN(2, vec3, normal)
IN(3, vec4, color)
IN(4, vec2, texcoord)

// Return Output
OUT(0, vec4, diffuse)
OUT(1, vec4, normal)
OUT(2, vec4, position)

// OBJECT(1, object, in_object)
SAMPLER_ARRAY(2, textures)

layout (scalar, binding = 1) readonly buffer uniform_objects { Object item[]; } _sb_objects;
Object object = _sb_objects.item[in_object];

void main() 
{
	uint texture0 = object.textures[TEX_SLOT_DIFFUSE];

	vec3 tint = mix(vec3(1), in_color.rgb, in_color.a);
	out_diffuse = vec4(texture_array(textures, texture0, in_texcoord).rgb * tint, 1);
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
