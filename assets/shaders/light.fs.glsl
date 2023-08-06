#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

//
// Lighting uniforms
//

layout (std430, binding = 1) readonly buffer LightBuffer {
	LightSettings settings;
	float[LIGHT_PADDING] _padding;
	Light item[];
} lights;

layout (binding = 2) uniform sampler2D tex_diffuse;
layout (binding = 3) uniform sampler2D tex_normal;
layout (binding = 4) uniform sampler2D tex_position;
layout (binding = 5) uniform sampler2D tex_occlusion;

// the variable-sized array must have the largest binding id
#define SHADOWMAP_SAMPLER shadowmaps
layout (binding = 6) uniform sampler2D[] shadowmaps;

#include "lib/lighting.glsl"

layout (location = 0) in vec2 v_texcoord0;

//
// Fragment output
//

layout (location = 0) out vec4 color;

void main() {
	// unpack data from geometry buffer
	vec3 viewPos = texture(tex_position, v_texcoord0).xyz;
	vec3 viewNormal = unpack_normal(texture(tex_normal, v_texcoord0).xyz);

	vec4 gcolor = texture(tex_diffuse, v_texcoord0);
	vec3 diffuseColor = gcolor.rgb;
	float occlusion = gcolor.a;

	vec3 position = getWorldPosition(viewPos);
	vec3 normal = getWorldNormal(viewNormal);

	float ssao = texture(tex_occlusion, v_texcoord0).r;

	// accumulate lighting
	vec3 lightColor = ambientLight(lights.settings, occlusion * ssao);
	int lightCount = lights.settings.Count;
	for(int i = 0; i < lightCount; i++) {
		lightColor += calculateLightColor(lights.item[i], position, normal, viewPos.z, lights.settings);
	}

	// linearize gbuffer diffuse
	vec3 linearDiffuse = pow(diffuseColor, vec3(2.2));

	// write shaded fragment color
	color = vec4(lightColor * linearDiffuse, 1);
}
