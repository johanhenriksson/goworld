#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

OBJECT(1, object)
SAMPLER_ARRAY(3, textures)

IN(4, vec2, texcoord0)
IN(5, flat uint, texture_id)

const int tile_count = 8;

void main() 
{
	uint texture0 = object.textures[0];

	// unpack tile coordinates
	vec2 t = vec2((in_texture_id + in_texcoord0.x) / tile_count, in_texcoord0.y);

	vec4 tint = vec4(in_color.rgb * in_color.a, 1);
	vec4 tex = vec4(texture_array(textures, texture0, t).rgb, 1);
	out_diffuse =tex * tint;

	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
