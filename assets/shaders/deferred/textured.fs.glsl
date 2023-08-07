#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/deferred_fragment.glsl"

STORAGE_BUFFER(1, Object, objects)
SAMPLER_ARRAY(3, textures)

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = objects.item[in_object].textures[0];

	out_diffuse = vec4(texture(textures[texture0], texcoord0).rgb, 1);
	out_normal = pack_normal(in_normal);
	out_position = vec4(in_position, 1);
}
