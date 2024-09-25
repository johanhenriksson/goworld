#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"
#include "lib/lighting.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)
OBJECT(1, object)
LIGHTS(2, lights)
SAMPLER_ARRAY(3, textures)

void main() 
{
	uint texture0 = object.textures[TEX_SLOT_DIFFUSE];
	vec4 albedo = texture_array(textures, texture0, in_texcoord) * in_color;

	// discard low alpha fragments
	if (albedo.a < 0.01) {
		discard;
	}

	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings, 1);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], in_world_position, in_world_normal, in_view_position.z, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor * lightColor, 1);
}
