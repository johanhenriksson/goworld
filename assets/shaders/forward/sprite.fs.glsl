#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/forward_fragment.glsl"

CAMERA(0, camera)
STORAGE_BUFFER(1, Object, objects)
SAMPLER_ARRAY(3, textures)

void main() 
{
	vec2 texcoord0 = in_color.xy;
	uint texture0 = objects.item[in_object].textures[0];
	vec4 albedo = texture(textures[texture0], texcoord0);

	if (albedo.a < 0.5) { 
		discard; 
	}

    // gamma correct & write fragment
	vec3 linearColor = pow(albedo.rgb, vec3(gamma));
    out_diffuse = vec4(linearColor, albedo.a);
}
