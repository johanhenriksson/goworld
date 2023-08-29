#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

STORAGE_BUFFER(1, Object, objects)
SAMPLER_ARRAY(3, textures)

IN(4, vec4, weights)

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = objects.item[in_object].textures[0];
	uint texture1 = objects.item[in_object].textures[1];
	uint texture2 = objects.item[in_object].textures[2];
	uint texture3 = objects.item[in_object].textures[3];

	vec3 color = in_weights.x * texture(textures[texture0], texcoord0).rgb +
		in_weights.y * texture(textures[texture1], texcoord0).rgb +
		in_weights.z * texture(textures[texture2], texcoord0).rgb +
		in_weights.w * texture(textures[texture3], texcoord0).rgb;

	color /= (in_weights.x + in_weights.y + in_weights.z + in_weights.w);

	out_diffuse = vec4(color, 1);
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
