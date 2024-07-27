#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/lighting.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)

float arc(vec3 a, vec3 b, float threshold) {
	return max(0, dot(a, b) - threshold) / (1 - threshold);
}

void main() 
{
	vec3 sky_angle = normalize(in_world_position - camera.Eye.xyz);

	vec3 sun_angle = normalize(vec3(0, 1, -1));
	vec3 sun_color = vec3(1, 0.8, 0.4);
	float sun_intensity = 16;
	float sun_halo = sun_intensity * 0.007;

	vec3 top = vec3(0.53, 0.69, 0.85);
	vec3 horizon = vec3(0.06, 0.18, 0.49);

	float frac = max(0, dot(sky_angle, vec3(0,1,0)));

	vec3 sky_color = mix(top, horizon, frac) + 
		asin(arc(sun_angle, sky_angle, 0.9985)) * sun_intensity * sun_color +
		asin(arc(sun_angle, sky_angle, 0.99)) * sun_halo * sun_color;

    out_diffuse = vec4(sky_color, in_color.a);
}

