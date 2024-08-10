#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

OBJECT(1, object)
SAMPLER_ARRAY(3, textures)

IN(4, vec4, weights)
IN(5, flat uvec4, indices)

// texture coordinate is stored in the color attribute.
vec2 texcoord = in_color.xy;

vec4 sample_texture(uint index, float weight) {
	if (weight <= 0) {
		return vec4(0);
	}
	uint texture_id = object.textures[index];
	vec4 color = texture_array(textures, texture_id, texcoord);
	return weight * vec4(color.rgb, 1);
}

void main() 
{
	vec4 color = vec4(0) +
		sample_texture(in_indices.x, in_weights.x) +
		sample_texture(in_indices.y, in_weights.y) +
		sample_texture(in_indices.z, in_weights.z) +
		sample_texture(in_indices.w, in_weights.w);

	// the sum of the weights are stored in the alpha channel.
	// divide by alpha to get the average color.
	out_diffuse = color / color.a;

	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
