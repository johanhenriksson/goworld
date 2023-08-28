#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/lighting.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)
LIGHT_BUFFER(2, lights)
SAMPLER_ARRAY(3, textures)

void main() 
{
	vec3 top = vec3(0.53, 0.69, 0.85);
	vec3 horizon = vec3(0.06, 0.18, 0.49);

	float frac = max(0, dot(in_world_normal, vec3(0,1,0)));

	vec3 sky_color = mix(top, horizon, frac);

    out_diffuse = vec4(sky_color, in_color.a);
}
