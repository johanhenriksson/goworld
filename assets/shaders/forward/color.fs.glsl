#version 450
#extension GL_GOOGLE_include_directive : enable

#include "../lib/common.glsl"
#include "../lib/material.glsl"
#include "../lib/fragment.glsl"

#define SHADOWMAP_SAMPLER Textures
#include "../lib/lighting.glsl"

// Varying
layout (location = 4) in vec3 wnormal;
layout (location = 5) in vec3 wposition;

void main() 
{
	int lightCount = 4;
	vec3 lightColor = vec3(0);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], wposition, wnormal, position0.z, 1);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(color0.rgb, vec3(gamma));
    diffuse = vec4(linearColor * lightColor, color0.a);

    // update gbuffer
    normal = pack_normal(normal0);
    position = vec4(position0, 1);
}
