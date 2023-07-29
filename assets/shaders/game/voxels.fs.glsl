#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable
#extension GL_GOOGLE_include_directive : enable

#include "../lib/uniforms.glsl"
#include "../lib/fragment.glsl"

void main() 
{
	diffuse = vec4(color0.rgb * color0.a, 1);

	vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
	normal = pack_normal;

	position = vec4(position0, 1);
}
