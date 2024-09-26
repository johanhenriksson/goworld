#extension GL_ARB_separate_shader_objects : require
#extension GL_ARB_shading_language_420pack : require
#extension GL_EXT_nonuniform_qualifier : require
#extension GL_EXT_buffer_reference : enable
#extension GL_EXT_buffer_reference2 : enable
#extension GL_EXT_shader_explicit_arithmetic_types : require
#extension GL_EXT_scalar_block_layout : require
#extension GL_EXT_shader_8bit_storage : require

#define MAX_TEXTURES 16

const float gamma = 2.2;

#define SAMPLER_ARRAY(idx,name) \
	layout (binding = idx) uniform sampler2D[] name; \
	float _shadow_texture(uint index, vec2 point) { return texture(name[nonuniformEXT(index)], point).r; } \
	vec2 _shadow_size(uint index) { return textureSize(name[index], 0).xy; }

#define texture_array(name,index,point) texture(name[nonuniformEXT(index)], point)

#define SAMPLER(idx,name) layout (binding = idx) uniform sampler2D tex_ ## name;

#define UNIFORM(idx,name,body) layout (binding = idx) uniform uniform_ ## name body name;

#define STORAGE_BUFFER(idx,type,name) layout (std430, binding = idx) readonly buffer uniform_ ## name { type item[]; } name;

#define CAMERA(idx,name) layout (binding = idx) uniform Camera { \
	mat4 Proj; \
	mat4 View; \
	mat4 ViewProj; \
	mat4 ProjInv; \
	mat4 ViewInv; \
	mat4 ViewProjInv; \
	vec4 Eye; \
	vec4 Forward; \
	vec2 Viewport; \
	float Delta; \
	float Time; \
} name;

#define IN(idx,type,name) layout (location = idx) in type in_ ## name;
#define OUT(idx,type,name) layout (location = idx) out type out_ ## name;

vec3 unpack_normal(vec3 packed_normal) {
	return normalize(2.0 * packed_normal - 1);
}

vec4 pack_normal(vec3 normal) {
	return vec4((normal + 1.0) / 2.0, 1);
}
