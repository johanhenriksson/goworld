#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/ui.glsl"

SAMPLER_ARRAY(2, textures)

IN(0, vec4, color)
IN(1, vec2, texcoord)
IN(2, vec2, position)
IN(3, flat vec2, center)
IN(4, flat vec2, half_size)
IN(5, flat uint, quad_index)
OUT(0, vec4, color)

void main() 
{
	Quad quad = quads.item[in_quad_index];

	// shrink the rectangle's half-size that is used for distance calculations 
	// otherwise the underlying primitive will cut off the falloff too early.
	vec2 softness_padding = vec2(max(0, quad.edge_softness*2-1),
								 max(0, quad.edge_softness*2-1));

	// sample distance to rect at position
	float dist = RoundedRectSDF(in_position,
								in_center,
								in_half_size-softness_padding,
								quad.corner_radius);

	// map distance to a blend factor
	float sdf_factor = 1 - smoothstep(0, 2*quad.edge_softness, dist);

	float border_factor = 1.f;
	if(quad.border > 0)
	{
		vec2 interior_half_size = in_half_size - vec2(quad.border);
		float interior_radius = quad.corner_radius - quad.border;

		// calculate sample distance from interior
		float inside_d = RoundedRectSDF(in_position,
										in_center,
										interior_half_size-
										softness_padding,
										interior_radius);

		// map distance => factor
		float inside_f = smoothstep(0, 2*quad.edge_softness, inside_d);
		border_factor = inside_f;
	}

	vec4 sample0 = texture_array(textures, quad.texture, in_texcoord);
	if(sample0.a == 0) {
		discard;
	}
	out_color = in_color * sample0 * sdf_factor * border_factor;
}

