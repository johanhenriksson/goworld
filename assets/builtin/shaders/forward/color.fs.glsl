#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/lighting.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)
LIGHTS(2, lights)
SAMPLER_ARRAY(3, textures)

void main() 
{
	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings, 1);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], in_world_position, in_world_normal, in_view_position.z, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(in_color.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor * lightColor, in_color.a);
}
