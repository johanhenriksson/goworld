#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"
#include "lib/lighting.glsl"

CAMERA(0, camera)
LIGHTS(1, lights)
SAMPLER(2, diffuse)
SAMPLER(3, normal)
SAMPLER(4, position)
SAMPLER(5, occlusion)
SAMPLER_ARRAY(6, shadowmaps)

IN(0, vec2, texcoord)
OUT(0, vec4, color)

vec3 getWorldPosition(vec3 viewPos); 
vec3 getWorldNormal(vec3 viewNormal);

void main() {
	// unpack data from geometry buffer
	vec3 viewPos = texture(tex_position, in_texcoord).xyz;
	vec3 viewNormal = unpack_normal(texture(tex_normal, in_texcoord).xyz);

	vec4 gcolor = texture(tex_diffuse, in_texcoord);
	vec3 diffuseColor = gcolor.rgb;
	float occlusion = gcolor.a;

	vec3 position = getWorldPosition(viewPos);
	vec3 normal = getWorldNormal(viewNormal);

	float ssao = texture(tex_occlusion, in_texcoord).r;
	if (ssao == 0) {
		ssao = 1;
	}

	// accumulate lighting
	vec3 lightColor = ambientLight(lights.settings, occlusion * ssao);
	int lightCount = lights.settings.Count;
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], position, normal, viewPos.z, lights.settings);
	}

	// linearize gbuffer diffuse
	vec3 linearDiffuse = pow(diffuseColor, vec3(2.2));

	// write shaded fragment color
	out_color = vec4(lightColor * linearDiffuse, 1);
}

vec3 getWorldPosition(vec3 viewPos) {
	vec4 pos_ws = camera.ViewInv * vec4(viewPos, 1);
	return pos_ws.xyz / pos_ws.w;
}

vec3 getWorldNormal(vec3 viewNormal) {
	vec4 worldNormal = camera.ViewInv * vec4(viewNormal, 0);
	return normalize(worldNormal.xyz);
}
