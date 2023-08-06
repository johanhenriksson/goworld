#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

const float gamma = 2.2;

#include "camera.glsl"

#define POINT_LIGHT 1
#define DIRECTIONAL_LIGHT 2

const int SHADOW_CASCADES = 4;
const float SHADOW_POWER = 60;

struct Attenuation {
	float Constant;
	float Linear;
	float Quadratic;
};

struct Light {
	mat4 ViewProj[SHADOW_CASCADES];  // 16 x 4 = 64 -> 64
	int Shadowmap[SHADOW_CASCADES];  // 4 x 4  = 16 -> 80
	float Distance[SHADOW_CASCADES]; // 4 x 4  = 16 -> 96

	vec4 Color;                      // 16 -> 112
	vec4 Position; 					 // 16 -> 128
	uint Type; 						 // 4  -> 132
	float Intensity; 				 // 4  -> 136
	float Range; 					 // 4  -> 140
	Attenuation Attenuation; 		 // 12 -> 152
};

struct LightSettings {
	vec4 AmbientColor;
	float AmbientIntensity;
	int Count;
	int ShadowSamples;
	float ShadowSampleRadius;
	float ShadowBias;
	float NormalOffset;
};

#define LIGHT_PADDING 76

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
