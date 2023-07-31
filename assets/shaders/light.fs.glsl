#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable
#extension GL_GOOGLE_include_directive : enable

#include "lib/shadows.glsl"

layout (std140, binding = 4) uniform Camera {
	mat4 Proj;
	mat4 View;
	mat4 ViewProj;
	mat4 ProjInv;
	mat4 ViewInv;
	mat4 ViewProjInv;
	vec3 Eye;
	vec3 Forward;
} camera;

layout (std430, binding = 5) readonly buffer LightBuffer {
	Light lights[];
} ssbo;


layout(push_constant) uniform constants
{
	int Count;
} push;

layout (input_attachment_index = 0, binding = 0) uniform subpassInput tex_diffuse;
layout (input_attachment_index = 1, binding = 1) uniform subpassInput tex_normal;
layout (input_attachment_index = 2, binding = 2) uniform subpassInput tex_position;
layout (input_attachment_index = 3, binding = 3) uniform subpassInput tex_depth;

layout (location = 0) out vec4 color;

const bool debug = false;

vec3 getWorldPosition(vec3 viewPos) {
	// transform view space to world space
	vec4 pos_ws = camera.ViewInv * vec4(viewPos, 1);
	return pos_ws.xyz / pos_ws.w;
}

float getDepth() {
	return subpassLoad(tex_position).z;
}

vec3 getWorldNormal() {
	// sample normal vector and transform it into world space
	vec3 viewNormal = normalize(2.0 * subpassLoad(tex_normal).rgb - 1); // normals [-1,1] 
	vec4 worldNormal = camera.ViewInv * vec4(viewNormal, 0);
	return normalize(worldNormal.xyz);
}


/* calculates lighting contribution from a point light source */
float calculatePointLightContrib(Light light, vec3 surfaceToLight, float distanceToLight, vec3 normal) {
	if (distanceToLight > light.Range) {
		return 0.0;
	}

	/* calculate normal coefficient */
	float normalCoef = max(0.0, dot(normal, surfaceToLight));

	/* light attenuation as a function of range and distance */
	float attenuation = light.Attenuation.Constant +
						light.Attenuation.Linear * distanceToLight +
						light.Attenuation.Quadratic * pow(distanceToLight, 2);
	attenuation = 1.0 / attenuation;

	/* multiply and return light contribution */
	return normalCoef * attenuation;
}

void main() {
	// unpack data from geometry buffer
	vec3 viewPos = subpassLoad(tex_position).xyz;
	vec4 t = subpassLoad(tex_diffuse);
	vec3 diffuseColor = t.rgb;
	float occlusion = t.a;

	vec3 position = getWorldPosition(viewPos);
	vec3 normal = getWorldNormal();

	vec3 lightColor = vec3(0);
	for(int i = 0; i < push.Count; i++) {
		Light light = ssbo.lights[i];

		// calculate contribution from the light source
		float contrib = 0.0;
		float shadow = 1.0;
		if (light.Type == AMBIENT_LIGHT) {
			contrib = 1;
		}
		else if (light.Type == DIRECTIONAL_LIGHT) {
			// directional lights store the direction in the position uniform
			// i.e. the light coming from the position, shining towards the origin
			vec3 lightDir = normalize(light.Position.xyz);
			vec3 surfaceToLight = -lightDir;
			contrib = max(dot(surfaceToLight, normal), 0.0);

			float bias = shadow_bias * max(0.0, 1.0 - dot(normal, lightDir));
			position += normal * normal_offset;
			shadow = blendCascades(light, position, viewPos.z, bias, 2, shadow_samples, sample_radius);

			if (debug) {
				int index = -1;
				for(int i = 0; i < SHADOW_CASCADES; i++) {
					if (viewPos.z < light.Distance[i]) {
						index = i;
						break;
					}
				}
				diffuseColor = mix(vec3(0,1,0), vec3(1,0,0), float(index) / (SHADOW_CASCADES - 1));
			}
		}
		else if (light.Type == POINT_LIGHT) {
			// calculate light vector & distance
			vec3 surfaceToLight = light.Position.xyz - position;
			float distanceToLight = length(surfaceToLight);
			surfaceToLight = normalize(surfaceToLight);
			contrib = calculatePointLightContrib(light, surfaceToLight, distanceToLight, normal);
		} 

		lightColor += light.Color.rgb * light.Intensity * contrib * shadow * occlusion;
	}

	// linearize gbuffer diffuse
	vec3 linearDiffuse = pow(diffuseColor, vec3(2.2));

	lightColor *= linearDiffuse;

	// lightColor *= mix(1, ssao, ssao_amount);

	// write fragment color & restore depth buffer
	color = vec4(lightColor, 1.0);
}
