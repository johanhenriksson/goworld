#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/lighting.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)
LIGHT_BUFFER(2, lights)
SAMPLER_ARRAY(3, textures)

// Varying
layout (location = 4) in vec3 wnormal;
layout (location = 5) in vec3 wposition;

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = objects.item[in_object].textures[0];
	vec4 albedo = texture(textures[texture0], texcoord0);

	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings, 1);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], wposition, wnormal, position0.z, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor * lightColor, albedo.a);
}
