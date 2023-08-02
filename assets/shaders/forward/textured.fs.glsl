#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/material.glsl"
#include "lib/fragment.glsl"

#define SHADOWMAP_SAMPLER Textures
#include "lib/lighting.glsl"

// Varying
layout (location = 4) in vec3 wnormal;
layout (location = 5) in vec3 wposition;

void main() 
{
	vec2 texcoord0 = color0.xy;
	uint texture0 = objects.item[objectIndex].textures[0];
	vec4 albedo = texture(Textures[texture0], texcoord0);

	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], wposition, wnormal, position0.z, 1, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    diffuse = vec4(linearColor * lightColor, albedo.a);

    // update gbuffer
    normal = pack_normal(normal0);
    position = vec4(position0, 1);
}
