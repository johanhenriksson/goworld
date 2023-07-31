#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

const float gamma = 2.2;

#include "camera.glsl"

vec3 unpack_normal(vec3 packed_normal) {
	return normalize(2.0 * packed_normal - 1);
}

vec4 pack_normal(vec3 normal) {
	return vec4((normal + 1.0) / 2.0, 1);
}

vec3 getWorldPosition(vec3 viewPos) {
	vec4 pos_ws = camera.ViewInv * vec4(viewPos, 1);
	return pos_ws.xyz / pos_ws.w;
}

vec3 getWorldNormal(vec3 viewNormal) {
	vec4 worldNormal = camera.ViewInv * vec4(viewNormal, 0);
	return normalize(worldNormal.xyz);
}
