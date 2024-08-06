#version 450
#extension GL_GOOGLE_include_directive : enable

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
	vec4 tex1 = texture(textures[object.textures[0]], texcoord0);
	vec4 tex2 = texture(textures[object.textures[1]], texcoord0);

	// fade between the two textures
	float alpha = sin(in_world_position.x * camera.Time * 0.11) * cos(in_world_position.z * camera.Time * 0.07) * 0.5 + 0.5;
	vec4 albedo = tex1 * alpha + tex2 * (1.0 - alpha);

	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings, 1);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], in_world_position, in_world_normal, in_view_position.z, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor * lightColor, albedo.a + 0.1);
}
