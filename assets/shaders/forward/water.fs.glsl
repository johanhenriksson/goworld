#version 450

#include "lib/common.glsl"
#include "lib/lighting.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)
OBJECT(1, object)
LIGHTS(2, lights)
SAMPLER_ARRAY(3, textures)

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = object.textures[0];
	vec4 albedo = texture_array(textures, texture0, texcoord0);

	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings, 1);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], in_world_position, in_world_normal, in_view_position.z, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor * lightColor, albedo.a+0.2);
}
