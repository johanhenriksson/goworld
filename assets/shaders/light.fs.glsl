#version 450
#extension GL_GOOGLE_include_directive : enable

#include "lib/common.glsl"

//
// Lighting uniforms
//

layout (std430, binding = 1) readonly buffer LightBuffer {
	Light item[];
} lights;

layout (input_attachment_index = 0, binding = 2) uniform subpassInput tex_diffuse;
layout (input_attachment_index = 1, binding = 3) uniform subpassInput tex_normal;
layout (input_attachment_index = 2, binding = 4) uniform subpassInput tex_position;

// the variable-sized array must have the largest binding id
#define SHADOWMAP_SAMPLER shadowmaps
layout (binding = 5) uniform sampler2D[] shadowmaps;

#include "lib/lighting.glsl"

//
// Push constants
//

// the shader expects the number of in-use lights as a push constant
layout(push_constant) uniform constants {
	int Count;
} push;

//
// Fragment output
//

layout (location = 0) out vec4 color;

void main() {
	// unpack data from geometry buffer
	vec3 viewPos = subpassLoad(tex_position).xyz;
	vec3 viewNormal = unpack_normal(subpassLoad(tex_normal).xyz);

	vec4 gcolor = subpassLoad(tex_diffuse);
	vec3 diffuseColor = gcolor.rgb;
	float occlusion = gcolor.a;

	vec3 position = getWorldPosition(viewPos);
	vec3 normal = getWorldNormal(viewNormal);

	// accumulate lighting
	vec3 lightColor = vec3(0);
	for(int i = 0; i < push.Count; i++) {
		lightColor += calculateLightColor(lights.item[i], position, normal, viewPos.z, occlusion);
	}

	// linearize gbuffer diffuse
	vec3 linearDiffuse = pow(diffuseColor, vec3(2.2));

	// write shaded fragment color
	color = vec4(lightColor * linearDiffuse, 1);
}
