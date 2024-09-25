#version 450

#include "lib/common.glsl"
#include "lib/objects.glsl"
#include "lib/lighting.glsl"

IN(0, flat uint, object)
IN(1, vec4, color)
IN(2, vec2, texcoord)
IN(3, vec3, view_position)
IN(4, vec3, world_normal)
IN(5, vec3, world_position)

// Return Output
OUT(0, vec4, diffuse)

CAMERA(0, camera)
OBJECT(1, object, in_object)
LIGHTS(2, lights)
SAMPLER_ARRAY(3, textures)

void main() 
{
	uint texture0 = object.textures[TEX_SLOT_DIFFUSE];
	vec4 albedo = texture_array(textures, texture0, in_texcoord) * in_color;

	if (albedo.a < 0.1) { 
		discard; 
	}

	int lightCount = lights.settings.Count;
	vec3 lightColor = ambientLight(lights.settings, 1);
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], in_world_position, in_world_normal, in_view_position.z, lights.settings);
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor * lightColor, albedo.a);
}
