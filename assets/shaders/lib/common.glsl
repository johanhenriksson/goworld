#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_EXT_nonuniform_qualifier : enable

vec4 pack_normal(vec3 normal) {
	return vec4((normal + 1.0) / 2.0, 1);
}
